package ex01_mcp_server

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"leetagent/common"
)

type echoHandler struct{}

func (echoHandler) Handle(ctx context.Context, args json.RawMessage) (string, error) {
	var p struct{ Text string }
	if err := json.Unmarshal(args, &p); err != nil {
		return "", err
	}
	if p.Text == "" {
		return "", errors.New("text is required")
	}
	return "echo: " + p.Text, nil
}

func run(t *testing.T, input string) []map[string]any {
	t.Helper()
	s := &Server{Name: "forum-mcp", Version: "0.1.0"}
	s.Register(common.NewFunctionTool("zz_last", "排序测试", `{"type":"object"}`), echoHandler{})
	s.Register(common.NewFunctionTool("echo", "原样返回 text", `{"type":"object","properties":{"text":{"type":"string"}}}`), echoHandler{})
	var out strings.Builder
	if err := s.Serve(context.Background(), strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	var res []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("response line is not JSON: %q", line)
		}
		res = append(res, m)
	}
	return res
}

func TestHandshakeAndList(t *testing.T) {
	in := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"t"}}}
{"jsonrpc":"2.0","method":"notifications/initialized"}

{"jsonrpc":"2.0","id":"two","method":"tools/list"}
`
	res := run(t, in)
	if len(res) != 2 {
		t.Fatalf("want 2 responses (notification gets none), got %d: %v", len(res), res)
	}
	init := res[0]["result"].(map[string]any)
	if res[0]["id"].(float64) != 1 || init["protocolVersion"] != ProtocolVersion {
		t.Fatalf("init = %v", res[0])
	}
	if init["serverInfo"].(map[string]any)["name"] != "forum-mcp" {
		t.Fatal("serverInfo.name")
	}
	if _, ok := init["capabilities"].(map[string]any)["tools"]; !ok {
		t.Fatal("capabilities.tools missing")
	}
	if res[1]["id"] != "two" {
		t.Fatal("string ids must be echoed back")
	}
	tools := res[1]["result"].(map[string]any)["tools"].([]any)
	first := tools[0].(map[string]any)
	if len(tools) != 2 || first["name"] != "echo" || first["inputSchema"].(map[string]any)["type"] != "object" {
		t.Fatalf("tools = %v", tools)
	}
}

func TestCallAndErrors(t *testing.T) {
	in := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"echo","arguments":{"text":"你好"}}}
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"echo","arguments":{}}}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"nope","arguments":{}}}
{"jsonrpc":"2.0","id":4,"method":"resources/list"}
{this is not json
{"jsonrpc":"2.0","id":5,"method":"ping"}
`
	res := run(t, in)
	if len(res) != 6 {
		t.Fatalf("got %d responses", len(res))
	}
	ok := res[0]["result"].(map[string]any)
	if ok["isError"] != false || ok["content"].([]any)[0].(map[string]any)["text"] != "echo: 你好" {
		t.Fatalf("call ok = %v", ok)
	}
	bad := res[1]["result"].(map[string]any)
	if bad["isError"] != true || bad["content"].([]any)[0].(map[string]any)["text"] != "text is required" {
		t.Fatalf("tool error should be a result with isError: %v", bad)
	}
	wantErr := []struct {
		idx  int
		code float64
		msg  string
	}{{2, CodeInvalidParams, "unknown tool: nope"}, {3, CodeMethodNotFound, "method not found: resources/list"}, {4, CodeParseError, "parse error"}}
	for _, w := range wantErr {
		e, _ := res[w.idx]["error"].(map[string]any)
		if e == nil || e["code"].(float64) != w.code || e["message"] != w.msg {
			t.Errorf("response %d = %v", w.idx, res[w.idx])
		}
	}
	if v, present := res[4]["id"]; !present || v != nil {
		t.Errorf("parse error must have id null: %v", res[4])
	}
	if _, ok := res[5]["result"]; !ok {
		t.Errorf("ping = %v", res[5])
	}
}
