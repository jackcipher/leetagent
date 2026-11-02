//go:build !solution

package ex02_mcp_client

import (
	"context"
	"encoding/json"
	"io"

	"leetagent/common"
)

// MCPTool 是 tools/list 返回的工具描述。
type MCPTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// RPCError 是服务端返回的 JSON-RPC 错误，实现 error：Error() 返回 "rpc error <code>: <message>"。
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	// TODO
	return ""
}

// ToolError 表示工具执行失败（result.isError == true），Error() 返回工具输出的文本。
type ToolError struct{ Text string }

func (e *ToolError) Error() string { return e.Text }

// Client 是 stdio MCP 客户端：向 w 写请求行，从 r 读响应行。请求 id 从 1 开始递增（数字）。
// 读响应时跳过：空行、没有 id 的消息（服务器通知）、id 不匹配的消息。
type Client struct {
	// TODO
}

func NewClient(r io.Reader, w io.Writer) *Client {
	// TODO
	return &Client{}
}

// Initialize 发送 initialize（params: protocolVersion "2025-06-18", capabilities {}, clientInfo {name:"leetagent", version:"1.0"}），
// 收到结果后发送通知 notifications/initialized（无 id），返回 serverInfo.name。
func (c *Client) Initialize(ctx context.Context) (string, error) {
	// TODO
	return "", nil
}

// ListTools 调用 tools/list。
func (c *Client) ListTools(ctx context.Context) ([]MCPTool, error) {
	// TODO
	return nil, nil
}

// CallTool 调用 tools/call，把 content 中所有 type=="text" 的 text 用 "\n" 连接返回。
// isError 为 true 时返回 ("", &ToolError{连接后的文本})。argsJSON 为空时发送 {}。
func (c *Client) CallTool(ctx context.Context, name, argsJSON string) (string, error) {
	// TODO
	return "", nil
}

// Call 让 Client 满足 common.ToolRunner，等价于 CallTool。
func (c *Client) Call(ctx context.Context, name, argsJSON string) (string, error) {
	return c.CallTool(ctx, name, argsJSON)
}

// ToCommonTools 把 MCP 工具转成 OpenAI function tools（inputSchema -> parameters），可直接交给模型。
func ToCommonTools(tools []MCPTool) []common.Tool {
	// TODO
	return nil
}
