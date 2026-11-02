//go:build !solution

package ex01_mcp_server

import (
	"context"
	"encoding/json"
	"io"

	"leetagent/common"
)

// ProtocolVersion 是本练习实现的 MCP 协议版本。
const ProtocolVersion = "2025-06-18"

// JSON-RPC 错误码。
const (
	CodeParseError     = -32700
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
)

// Request / Response 是 JSON-RPC 2.0 消息。ID 为 nil 表示通知（notification），不需要回复。
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

// Handler 执行一个工具。
type Handler interface {
	Handle(ctx context.Context, args json.RawMessage) (string, error)
}

// Server 是一个通过 stdio（按行分隔的 JSON）通信的 MCP 服务器。
type Server struct {
	Name    string
	Version string
	// TODO: 工具注册表字段自行添加
}

// Register 注册工具（同名覆盖）。
func (s *Server) Register(def common.Tool, h Handler) {
	// TODO
}

// Serve 从 r 逐行读取请求，向 w 逐行写响应（每个响应是一行 JSON + "\n"），直到 r 结束（返回 nil）或 ctx 取消。
// 方法：
//
//	initialize      -> {"protocolVersion": ProtocolVersion, "capabilities": {"tools": {}}, "serverInfo": {"name","version"}}
//	tools/list      -> {"tools": [{"name","description","inputSchema"}...]}，按 name 排序
//	tools/call      params {"name","arguments"} ->
//	                   成功 {"content":[{"type":"text","text": 输出}], "isError": false}
//	                   Handler 出错 {"content":[{"type":"text","text": err.Error()}], "isError": true}
//	                   未知工具 -> error {CodeInvalidParams, "unknown tool: <name>"}
//	ping            -> {}
//	其它方法        -> error {CodeMethodNotFound, "method not found: <method>"}（通知则忽略）
//	无法解析的行    -> error {CodeParseError, "parse error"}，id 为 null
//	空行跳过；通知（没有 id）一律不回复。
func (s *Server) Serve(ctx context.Context, r io.Reader, w io.Writer) error {
	// TODO
	return nil
}
