// Package mcp is a small, dependency-free Model Context Protocol
// server: tools described by Go functions, served over Streamable HTTP
// or stdio. Giode uses it to expose a running app to AI agents (see
// Config.MCP in the root package), and it can be used on its own:
//
//	srv := mcp.NewServer("my-app", "1.0.0")
//	srv.AddTool(mcp.Func("greet", "Greets someone.",
//		func(ctx context.Context, in struct {
//			Name string `json:"name" description:"Who to greet"`
//		}) (any, error) {
//			return "Hello, " + in.Name, nil
//		}))
//	http.Handle("/mcp", srv)
//
// Only the tools capability is implemented; resources and prompts are
// advertised as empty lists.
package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

// LatestProtocolVersion is the newest MCP revision this server speaks.
const LatestProtocolVersion = "2025-06-18"

// supportedVersions are the revisions negotiated with clients, newest
// first.
var supportedVersions = []string{LatestProtocolVersion, "2025-03-26", "2024-11-05"}

// Tool is a callable tool.
type Tool struct {
	// Name identifies the tool; letters, digits, '_' and '-'.
	Name string
	// Description tells the model what the tool does and when to use
	// it.
	Description string
	// InputSchema is the JSON Schema of the arguments object. Nil
	// means no arguments.
	InputSchema map[string]any
	// ReadOnly hints that the tool doesn't modify anything.
	ReadOnly bool
	// Handler runs the tool with the raw arguments object.
	Handler func(ctx context.Context, args json.RawMessage) (*Result, error)
}

// Result is a tool call's result.
type Result struct {
	Content []Content `json:"content"`
	// IsError reports a tool-level failure the model should see (as
	// opposed to a protocol error).
	IsError bool `json:"isError,omitempty"`
}

// Content is one item of a tool result.
type Content struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// Text returns a result holding s.
func Text(s string) *Result {
	return &Result{Content: []Content{{Type: "text", Text: s}}}
}

// JSON returns a result holding v encoded as indented JSON.
func JSON(v any) (*Result, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return Text(string(b)), nil
}

// Image returns a result holding an image.
func Image(data []byte, mimeType string) *Result {
	return &Result{Content: []Content{{Type: "image", Data: base64.StdEncoding.EncodeToString(data), MimeType: mimeType}}}
}

// Errorf returns a tool-level error result.
func Errorf(format string, args ...any) *Result {
	r := Text(fmt.Sprintf(format, args...))
	r.IsError = true
	return r
}

// Func builds a Tool from a typed function. The input struct's exported
// fields become the arguments: their JSON names come from json tags,
// a field is required unless tagged omitempty, and description and
// enum ("a,b,c") tags document it. The output is returned as text if
// it's a string, as-is if it's a *Result, and as JSON otherwise.
func Func[In any](name, description string, fn func(ctx context.Context, in In) (any, error)) Tool {
	var zero In
	return Tool{
		Name:        name,
		Description: description,
		InputSchema: SchemaFor(reflect.TypeOf(zero)),
		Handler: func(ctx context.Context, args json.RawMessage) (*Result, error) {
			var in In
			if len(args) > 0 && string(args) != "null" {
				if err := json.Unmarshal(args, &in); err != nil {
					return Errorf("invalid arguments: %v", err), nil
				}
			}
			out, err := fn(ctx, in)
			if err != nil {
				return Errorf("%v", err), nil
			}
			switch v := out.(type) {
			case *Result:
				return v, nil
			case string:
				return Text(v), nil
			case nil:
				return Text("ok"), nil
			default:
				return JSON(v)
			}
		},
	}
}

// Server is an MCP server. Its zero value is not usable; create one
// with NewServer. A Server is safe for concurrent use.
type Server struct {
	name, version string
	// Instructions, when set, are sent to clients on initialize to
	// explain how to use the server.
	Instructions string
	// Token, when set, must be sent by HTTP clients as
	// "Authorization: Bearer <token>".
	Token string
	// AllowedOrigins lists browser origins allowed to call the HTTP
	// endpoint, besides localhost ones. Requests without an Origin
	// header (non-browser clients) are always allowed.
	AllowedOrigins []string

	mu    sync.RWMutex
	tools map[string]Tool
}

// NewServer returns a server identifying itself as name and version.
func NewServer(name, version string) *Server {
	return &Server{name: name, version: version, tools: map[string]Tool{}}
}

// AddTool registers (or replaces) a tool.
func (s *Server) AddTool(t Tool) {
	if t.InputSchema == nil {
		t.InputSchema = map[string]any{"type": "object", "properties": map[string]any{}}
	}
	s.mu.Lock()
	s.tools[t.Name] = t
	s.mu.Unlock()
}

// RemoveTool unregisters a tool.
func (s *Server) RemoveTool(name string) {
	s.mu.Lock()
	delete(s.tools, name)
	s.mu.Unlock()
}

// Tools returns the registered tools, sorted by name.
func (s *Server) Tools() []Tool {
	s.mu.RLock()
	out := make([]Tool, 0, len(s.tools))
	for _, t := range s.tools {
		out = append(out, t)
	}
	s.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// --- JSON-RPC -----------------------------------------------------------------

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSON-RPC error codes.
const (
	codeParse          = -32700
	codeInvalidRequest = -32600
	codeMethodNotFound = -32601
	codeInvalidParams  = -32602
	codeInternal       = -32603
)

// Handle processes one JSON-RPC message and returns the response to
// send, or nil for notifications.
func (s *Server) Handle(ctx context.Context, msg []byte) []byte {
	var req request
	if err := json.Unmarshal(msg, &req); err != nil {
		return encode(response{JSONRPC: "2.0", ID: json.RawMessage("null"), Error: &rpcError{codeParse, "parse error: " + err.Error()}})
	}
	if req.JSONRPC != "2.0" || req.Method == "" {
		if len(req.ID) == 0 {
			return nil
		}
		return encode(response{JSONRPC: "2.0", ID: req.ID, Error: &rpcError{codeInvalidRequest, "invalid request"}})
	}
	result, rerr := s.dispatch(ctx, req)
	if len(req.ID) == 0 {
		return nil // notification: no response, even on error
	}
	resp := response{JSONRPC: "2.0", ID: req.ID}
	if rerr != nil {
		resp.Error = rerr
	} else {
		resp.Result = result
	}
	return encode(resp)
}

func encode(r response) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		b, _ = json.Marshal(response{JSONRPC: "2.0", ID: r.ID, Error: &rpcError{codeInternal, err.Error()}})
	}
	return b
}

func (s *Server) dispatch(ctx context.Context, req request) (any, *rpcError) {
	switch req.Method {
	case "initialize":
		var p struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		_ = json.Unmarshal(req.Params, &p)
		version := LatestProtocolVersion
		for _, v := range supportedVersions {
			if v == p.ProtocolVersion {
				version = v
			}
		}
		res := map[string]any{
			"protocolVersion": version,
			"capabilities": map[string]any{
				"tools":     map[string]any{"listChanged": false},
				"resources": map[string]any{},
				"prompts":   map[string]any{},
			},
			"serverInfo": map[string]any{"name": s.name, "version": s.version},
		}
		if s.Instructions != "" {
			res["instructions"] = s.Instructions
		}
		return res, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		tools := s.Tools()
		list := make([]map[string]any, 0, len(tools))
		for _, t := range tools {
			entry := map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"inputSchema": t.InputSchema,
			}
			if t.ReadOnly {
				entry["annotations"] = map[string]any{"readOnlyHint": true}
			}
			list = append(list, entry)
		}
		return map[string]any{"tools": list}, nil
	case "tools/call":
		var p struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return nil, &rpcError{codeInvalidParams, err.Error()}
		}
		s.mu.RLock()
		t, ok := s.tools[p.Name]
		s.mu.RUnlock()
		if !ok {
			return nil, &rpcError{codeInvalidParams, "unknown tool " + p.Name}
		}
		res, err := callTool(ctx, t, p.Arguments)
		if err != nil {
			return Errorf("%v", err), nil
		}
		return res, nil
	case "resources/list":
		return map[string]any{"resources": []any{}}, nil
	case "resources/templates/list":
		return map[string]any{"resourceTemplates": []any{}}, nil
	case "prompts/list":
		return map[string]any{"prompts": []any{}}, nil
	case "logging/setLevel":
		return map[string]any{}, nil
	}
	if strings.HasPrefix(req.Method, "notifications/") {
		return nil, nil
	}
	return nil, &rpcError{codeMethodNotFound, "method not found: " + req.Method}
}

// callTool runs a handler, turning a panic into a tool error so one
// faulty tool can't take the app down.
func callTool(ctx context.Context, t Tool, args json.RawMessage) (res *Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tool %s panicked: %v", t.Name, r)
		}
	}()
	res, err = t.Handler(ctx, args)
	if err == nil && res == nil {
		res = Text("ok")
	}
	return res, err
}
