// Package solution 是 ex01_mcp_server 的参考答案。
package solution

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"sort"
	"strings"

	"leetagent/common"
)

const ProtocolVersion = "2025-06-18"

const (
	CodeParseError     = -32700
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type Handler interface {
	Handle(ctx context.Context, args json.RawMessage) (string, error)
}

type entry struct {
	def common.Tool
	h   Handler
}

type Server struct {
	Name    string
	Version string
	tools   map[string]entry
}

func (s *Server) Register(def common.Tool, h Handler) {
	if s.tools == nil {
		s.tools = map[string]entry{}
	}
	s.tools[def.Function.Name] = entry{def: def, h: h}
}

type toolInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type textContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type callResult struct {
	Content []textContent `json:"content"`
	IsError bool          `json:"isError"`
}

type callParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func (s *Server) Serve(ctx context.Context, r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	enc := json.NewEncoder(w)
	for sc.Scan() {
		if err := ctx.Err(); err != nil {
			return err
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			resp := Response{JSONRPC: "2.0", ID: json.RawMessage("null"), Error: &RPCError{CodeParseError, "parse error"}}
			if err := enc.Encode(resp); err != nil {
				return err
			}
			continue
		}
		if len(req.ID) == 0 || string(req.ID) == "null" {
			continue // notification
		}
		resp := s.dispatch(ctx, req)
		if err := enc.Encode(resp); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Server) dispatch(ctx context.Context, req Request) Response {
	resp := Response{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": ProtocolVersion,
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]string{"name": s.Name, "version": s.Version},
		}
	case "ping":
		resp.Result = map[string]any{}
	case "tools/list":
		resp.Result = map[string]any{"tools": s.list()}
	case "tools/call":
		resp.Result, resp.Error = s.call(ctx, req.Params)
	default:
		resp.Error = &RPCError{CodeMethodNotFound, "method not found: " + req.Method}
	}
	return resp
}

func (s *Server) list() []toolInfo {
	names := make([]string, 0, len(s.tools))
	for n := range s.tools {
		names = append(names, n)
	}
	sort.Strings(names)
	out := make([]toolInfo, 0, len(names))
	for _, n := range names {
		f := s.tools[n].def.Function
		out = append(out, toolInfo{Name: f.Name, Description: f.Description, InputSchema: f.Parameters})
	}
	return out
}

func (s *Server) call(ctx context.Context, raw json.RawMessage) (any, *RPCError) {
	var p callParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, &RPCError{CodeInvalidParams, "invalid params"}
	}
	e, ok := s.tools[p.Name]
	if !ok {
		return nil, &RPCError{CodeInvalidParams, "unknown tool: " + p.Name}
	}
	out, err := e.h.Handle(ctx, p.Arguments)
	if err != nil {
		return callResult{Content: []textContent{{"text", err.Error()}}, IsError: true}, nil
	}
	return callResult{Content: []textContent{{"text", out}}}, nil
}
