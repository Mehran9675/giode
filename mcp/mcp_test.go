package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type greetIn struct {
	Name  string `json:"name" description:"Who to greet"`
	Shout bool   `json:"shout,omitempty"`
	Mood  string `json:"mood,omitempty" enum:"happy,sad"`
}

func testServer() *Server {
	s := NewServer("test", "1.0.0")
	s.Instructions = "be nice"
	s.AddTool(Func("greet", "Greets someone.", func(ctx context.Context, in greetIn) (any, error) {
		if in.Name == "" {
			return nil, errors.New("name is required")
		}
		msg := "Hello, " + in.Name
		if in.Shout {
			msg = strings.ToUpper(msg)
		}
		return msg, nil
	}))
	s.AddTool(Func("boom", "Panics.", func(ctx context.Context, in struct{}) (any, error) {
		panic("kaboom")
	}))
	s.AddTool(Func("data", "Returns JSON.", func(ctx context.Context, in struct{}) (any, error) {
		return map[string]int{"n": 3}, nil
	}))
	return s
}

func call(t *testing.T, s *Server, msg string) map[string]any {
	t.Helper()
	out := s.Handle(context.Background(), []byte(msg))
	if out == nil {
		t.Fatalf("no response to %s", msg)
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("bad response %s: %v", out, err)
	}
	return m
}

func TestInitializeNegotiatesVersion(t *testing.T) {
	s := testServer()
	m := call(t, s, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"c","version":"1"}}}`)
	res := m["result"].(map[string]any)
	if res["protocolVersion"] != "2025-03-26" {
		t.Errorf("protocolVersion = %v, want the client's supported version", res["protocolVersion"])
	}
	if res["instructions"] != "be nice" {
		t.Errorf("instructions = %v", res["instructions"])
	}
	m = call(t, s, `{"jsonrpc":"2.0","id":2,"method":"initialize","params":{"protocolVersion":"1999-01-01"}}`)
	if v := m["result"].(map[string]any)["protocolVersion"]; v != LatestProtocolVersion {
		t.Errorf("unsupported version negotiated to %v, want %v", v, LatestProtocolVersion)
	}
}

func TestToolsListSchema(t *testing.T) {
	m := call(t, testServer(), `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	tools := m["result"].(map[string]any)["tools"].([]any)
	if len(tools) != 3 {
		t.Fatalf("tools = %d, want 3", len(tools))
	}
	var greet map[string]any
	for _, tl := range tools {
		if tl.(map[string]any)["name"] == "greet" {
			greet = tl.(map[string]any)
		}
	}
	schema := greet["inputSchema"].(map[string]any)
	req := schema["required"].([]any)
	if len(req) != 1 || req[0] != "name" {
		t.Errorf("required = %v, want [name]", req)
	}
	props := schema["properties"].(map[string]any)
	if props["name"].(map[string]any)["description"] != "Who to greet" {
		t.Errorf("description tag not applied: %v", props["name"])
	}
	if len(props["mood"].(map[string]any)["enum"].([]any)) != 2 {
		t.Errorf("enum tag not applied: %v", props["mood"])
	}
}

func TestToolsCall(t *testing.T) {
	s := testServer()
	text := func(m map[string]any) (string, bool) {
		res := m["result"].(map[string]any)
		isErr, _ := res["isError"].(bool)
		return res["content"].([]any)[0].(map[string]any)["text"].(string), isErr
	}
	got, isErr := text(call(t, s, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"greet","arguments":{"name":"Ada","shout":true}}}`))
	if got != "HELLO, ADA" || isErr {
		t.Errorf("greet = %q (error %v)", got, isErr)
	}
	if _, isErr := text(call(t, s, `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"greet","arguments":{}}}`)); !isErr {
		t.Error("handler error not reported as isError")
	}
	if got, isErr := text(call(t, s, `{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"boom"}}`)); !isErr || !strings.Contains(got, "kaboom") {
		t.Errorf("panic = %q (error %v), want a tool error", got, isErr)
	}
	if got, _ := text(call(t, s, `{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"data"}}`)); !strings.Contains(got, `"n": 3`) {
		t.Errorf("json result = %q", got)
	}
	m := call(t, s, `{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"nope"}}`)
	if m["error"] == nil {
		t.Error("unknown tool: want a protocol error")
	}
}

func TestNotificationsAndErrors(t *testing.T) {
	s := testServer()
	if out := s.Handle(context.Background(), []byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}`)); out != nil {
		t.Errorf("notification got a response: %s", out)
	}
	m := call(t, s, `{"jsonrpc":"2.0","id":1,"method":"does/not/exist"}`)
	if code := m["error"].(map[string]any)["code"].(float64); code != codeMethodNotFound {
		t.Errorf("code = %v, want %d", code, codeMethodNotFound)
	}
	m = call(t, s, `{not json`)
	if code := m["error"].(map[string]any)["code"].(float64); code != codeParse {
		t.Errorf("code = %v, want %d", code, codeParse)
	}
}

func TestHTTPTransport(t *testing.T) {
	s := testServer()
	s.Token = "secret"
	ts := httptest.NewServer(s)
	defer ts.Close()

	post := func(body, origin, token string) *http.Response {
		req, _ := http.NewRequest(http.MethodPost, ts.URL, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		return resp
	}
	init := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`
	if r := post(init, "", ""); r.StatusCode != http.StatusUnauthorized {
		t.Errorf("no token: status %d, want 401", r.StatusCode)
	}
	if r := post(init, "https://evil.example", "secret"); r.StatusCode != http.StatusForbidden {
		t.Errorf("foreign origin: status %d, want 403", r.StatusCode)
	}
	r := post(init, "http://localhost:3000", "secret")
	if r.StatusCode != http.StatusOK || r.Header.Get("Mcp-Session-Id") == "" {
		t.Errorf("initialize: status %d, session %q", r.StatusCode, r.Header.Get("Mcp-Session-Id"))
	}
	if r := post(`{"jsonrpc":"2.0","method":"notifications/initialized"}`, "", "secret"); r.StatusCode != http.StatusAccepted {
		t.Errorf("notification: status %d, want 202", r.StatusCode)
	}
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer secret")
	if resp, _ := http.DefaultClient.Do(req); resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("GET: status %d, want 405", resp.StatusCode)
	}
}

func TestStdioTransport(t *testing.T) {
	s := testServer()
	in := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping"}` + "\n\n" +
		`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n")
	var out bytes.Buffer
	if err := s.ServeStdio(context.Background(), in, &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("responses = %d (%q), want 2", len(lines), out.String())
	}
}
