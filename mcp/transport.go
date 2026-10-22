package mcp

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// maxMessageSize bounds a single JSON-RPC message.
const maxMessageSize = 8 << 20

// ServeHTTP implements the Streamable HTTP transport: each POST carries
// one JSON-RPC message and gets its response as application/json.
// Server-initiated streams (GET) aren't offered, which the
// specification allows.
//
// Requests from browsers are only accepted from localhost origins and
// AllowedOrigins, guarding a local server against DNS rebinding; when
// Token is set, every request must carry it as a bearer token.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.originAllowed(r.Header.Get("Origin")) {
		http.Error(w, "origin not allowed", http.StatusForbidden)
		return
	}
	if s.Token != "" {
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if subtle.ConstantTimeCompare([]byte(got), []byte(s.Token)) != 1 {
			w.Header().Set("WWW-Authenticate", `Bearer realm="mcp"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
	switch r.Method {
	case http.MethodPost:
	case http.MethodGet:
		w.Header().Set("Allow", "POST, DELETE")
		http.Error(w, "server-initiated streams are not supported", http.StatusMethodNotAllowed)
		return
	case http.MethodDelete:
		// Sessions hold no server state; ending one is a no-op.
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.Header().Set("Allow", "POST, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxMessageSize+1))
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}
	if len(body) > maxMessageSize {
		http.Error(w, "message too large", http.StatusRequestEntityTooLarge)
		return
	}
	body = bytes.TrimSpace(body)
	if len(body) > 0 && body[0] == '[' {
		http.Error(w, "batched messages are not supported", http.StatusBadRequest)
		return
	}
	resp := s.Handle(r.Context(), body)
	if isInitialize(body) {
		w.Header().Set("Mcp-Session-Id", newSessionID())
	}
	if resp == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func isInitialize(body []byte) bool {
	var m struct {
		Method string `json:"method"`
	}
	return json.Unmarshal(body, &m) == nil && m.Method == "initialize"
}

func newSessionID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func (s *Server) originAllowed(origin string) bool {
	if origin == "" {
		return true
	}
	for _, o := range s.AllowedOrigins {
		if strings.EqualFold(o, origin) {
			return true
		}
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := u.Hostname()
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// ListenAndServe serves the HTTP transport at addr under the path
// "/mcp" until ctx is canceled. It returns the address actually bound
// (useful with port 0) through ready, if non-nil, once listening.
func (s *Server) ListenAndServe(ctx context.Context, addr string, ready func(addr string)) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/mcp", s)
	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	if ready != nil {
		ready(ln.Addr().String())
	}
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()
	select {
	case <-ctx.Done():
		sctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(sctx)
		return nil
	case err := <-errc:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

// ServeStdio serves newline-delimited JSON-RPC messages from r,
// writing responses to w, until r is exhausted or ctx is canceled -
// the stdio transport, for clients that launch the app as a
// subprocess. Requests are handled concurrently, so a slow tool
// doesn't block others.
func (s *Server) ServeStdio(ctx context.Context, r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), maxMessageSize)
	var (
		wmu sync.Mutex
		wg  sync.WaitGroup
	)
	defer wg.Wait()
	for sc.Scan() {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		msg := append([]byte(nil), line...)
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp := s.Handle(ctx, msg)
			if resp == nil {
				return
			}
			wmu.Lock()
			defer wmu.Unlock()
			_, _ = w.Write(append(resp, '\n'))
		}()
	}
	return sc.Err()
}
