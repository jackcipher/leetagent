// Package solution 是 ex02_mcp_client 的参考答案。
package solution

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"leetagent/common"
)

type MCPTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string { return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message) }

type ToolError struct{ Text string }

func (e *ToolError) Error() string { return e.Text }

type Client struct {
	mu     sync.Mutex
	r      *bufio.Reader
	w      io.Writer
	nextID int
}

func NewClient(r io.Reader, w io.Writer) *Client {
	return &Client{r: bufio.NewReader(r), w: w}
}

type request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      *int   `json:"id,omitempty"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type response struct {
	ID     json.RawMessage `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *RPCError       `json:"error"`
}

func (c *Client) send(req request) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = c.w.Write(append(b, '\n'))
	return err
}

func (c *Client) rpc(ctx context.Context, method string, params any, out any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	c.nextID++
	id := c.nextID
	if err := c.send(request{JSONRPC: "2.0", ID: &id, Method: method, Params: params}); err != nil {
		return err
	}
	want := strconv.Itoa(id)
	for {
		line, err := c.r.ReadString('\n')
		if strings.TrimSpace(line) == "" {
			if err != nil {
				return err
			}
			continue
		}
		var resp response
		if json.Unmarshal([]byte(line), &resp) != nil || string(resp.ID) != want {
			if err != nil {
				return err
			}
			continue
		}
		if resp.Error != nil {
			return resp.Error
		}
		if out == nil {
			return nil
		}
		return json.Unmarshal(resp.Result, out)
	}
}

func (c *Client) Initialize(ctx context.Context) (string, error) {
	params := map[string]any{
		"protocolVersion": "2025-06-18",
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]string{"name": "leetagent", "version": "1.0"},
	}
	var res struct {
		ServerInfo struct {
			Name string `json:"name"`
		} `json:"serverInfo"`
	}
	if err := c.rpc(ctx, "initialize", params, &res); err != nil {
		return "", err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.send(request{JSONRPC: "2.0", Method: "notifications/initialized"}); err != nil {
		return "", err
	}
	return res.ServerInfo.Name, nil
}

func (c *Client) ListTools(ctx context.Context) ([]MCPTool, error) {
	var res struct {
		Tools []MCPTool `json:"tools"`
	}
	if err := c.rpc(ctx, "tools/list", map[string]any{}, &res); err != nil {
		return nil, err
	}
	return res.Tools, nil
}

func (c *Client) CallTool(ctx context.Context, name, argsJSON string) (string, error) {
	if strings.TrimSpace(argsJSON) == "" {
		argsJSON = "{}"
	}
	params := map[string]any{"name": name, "arguments": json.RawMessage(argsJSON)}
	var res struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		IsError bool `json:"isError"`
	}
	if err := c.rpc(ctx, "tools/call", params, &res); err != nil {
		return "", err
	}
	var texts []string
	for _, ct := range res.Content {
		if ct.Type == "text" {
			texts = append(texts, ct.Text)
		}
	}
	text := strings.Join(texts, "\n")
	if res.IsError {
		return "", &ToolError{Text: text}
	}
	return text, nil
}

func (c *Client) Call(ctx context.Context, name, argsJSON string) (string, error) {
	return c.CallTool(ctx, name, argsJSON)
}

func ToCommonTools(tools []MCPTool) []common.Tool {
	out := make([]common.Tool, len(tools))
	for i, t := range tools {
		out[i] = common.NewFunctionTool(t.Name, t.Description, string(t.InputSchema))
	}
	return out
}
