package ex02_mcp_client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"leetagent/common"
)

// fakeServer 读取客户端写来的每一行，记录下来，并按 method 给出回复。
type fakeServer struct {
	got []map[string]any
}

func (f *fakeServer) serve(r io.Reader, w io.Writer) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		var m map[string]any
		if json.Unmarshal(sc.Bytes(), &m) != nil {
			continue
		}
		f.got = append(f.got, m)
		id, ok := m["id"]
		if !ok {
			continue
		}
		idJSON, _ := json.Marshal(id)
		var result string
		switch m["method"] {
		case "initialize":
			// 先发一条服务器通知和一个空行，客户端应跳过
			io.WriteString(w, `{"jsonrpc":"2.0","method":"notifications/message","params":{"level":"info"}}`+"\n\n")
			// 再发一条 id 不匹配的，也应跳过
			io.WriteString(w, `{"jsonrpc":"2.0","id":999,"result":{}}`+"\n")
			result = `{"protocolVersion":"2025-06-18","capabilities":{"tools":{}},"serverInfo":{"name":"forum-mcp","version":"0.1"}}`
		case "tools/list":
			result = `{"tools":[{"name":"get_post","description":"读帖子","inputSchema":{"type":"object","properties":{"tid":{"type":"integer"}}}}]}`
		case "tools/call":
			params := m["params"].(map[string]any)
			switch params["name"] {
			case "get_post":
				result = `{"content":[{"type":"text","text":"标题: GMP"},{"type":"image","data":"xx"},{"type":"text","text":"正文..."}],"isError":false}`
			case "broken":
				result = `{"content":[{"type":"text","text":"forum returned 500"}],"isError":true}`
			default:
				io.WriteString(w, `{"jsonrpc":"2.0","id":`+string(idJSON)+`,"error":{"code":-32602,"message":"unknown tool: nope"}}`+"\n")
				continue
			}
		}
		io.WriteString(w, `{"jsonrpc":"2.0","id":`+string(idJSON)+`,"result":`+result+"}\n")
	}
}

func TestClientOverPipes(t *testing.T) {
	c2sR, c2sW := io.Pipe()
	s2cR, s2cW := io.Pipe()
	srv := &fakeServer{}
	go srv.serve(c2sR, s2cW)
	defer c2sW.Close()

	c := NewClient(s2cR, c2sW)
	var _ common.ToolRunner = c
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()
		s2cW.CloseWithError(errors.New("test timeout"))
	}()

	name, err := c.Initialize(ctx)
	if err != nil || name != "forum-mcp" {
		t.Fatalf("Initialize = %q, %v", name, err)
	}
	tools, err := c.ListTools(ctx)
	if err != nil || len(tools) != 1 || tools[0].Name != "get_post" {
		t.Fatalf("ListTools = %+v, %v", tools, err)
	}
	ct := ToCommonTools(tools)
	if ct[0].Type != "function" || ct[0].Function.Name != "get_post" || !strings.Contains(string(ct[0].Function.Parameters), `"tid"`) {
		t.Fatalf("ToCommonTools = %+v", ct)
	}
	out, err := c.CallTool(ctx, "get_post", `{"tid":1}`)
	if err != nil || out != "标题: GMP\n正文..." {
		t.Fatalf("CallTool = %q, %v", out, err)
	}
	_, err = c.Call(ctx, "broken", "")
	var te *ToolError
	if !errors.As(err, &te) || te.Text != "forum returned 500" {
		t.Fatalf("tool error = %v", err)
	}
	_, err = c.CallTool(ctx, "nope", "{}")
	var re *RPCError
	if !errors.As(err, &re) || re.Code != -32602 || err.Error() != "rpc error -32602: unknown tool: nope" {
		t.Fatalf("rpc error = %v", err)
	}

	// 检查客户端发出的消息
	if len(srv.got) < 6 {
		t.Fatalf("server saw %d messages", len(srv.got))
	}
	if srv.got[0]["id"].(float64) != 1 || srv.got[0]["method"] != "initialize" {
		t.Fatalf("first = %v", srv.got[0])
	}
	if pv := srv.got[0]["params"].(map[string]any)["protocolVersion"]; pv != "2025-06-18" {
		t.Fatalf("protocolVersion = %v", pv)
	}
	if _, hasID := srv.got[1]["id"]; hasID || srv.got[1]["method"] != "notifications/initialized" {
		t.Fatalf("second must be initialized notification: %v", srv.got[1])
	}
	if srv.got[2]["id"].(float64) != 2 {
		t.Fatalf("ids must increase: %v", srv.got[2])
	}
	args := srv.got[4]["params"].(map[string]any)["arguments"]
	if m, ok := args.(map[string]any); !ok || len(m) != 0 {
		t.Fatalf("empty args must be sent as {}: %v", args)
	}
	cancel()
	<-done
}
