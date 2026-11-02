// Package common 提供整个课程共用的最小类型、LLM 流式客户端和离线 mock LLM 服务器。
// 类型和 典型 Agent 项目的 llm 包 保持一致，这样练习写完可以直接搬回项目。
package common

import (
	"context"
	"encoding/json"
)

// Message 是 OpenAI 兼容的对话消息。
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// Tool 是请求里 tools[] 的一项。
type Tool struct {
	Type     string   `json:"type"`
	Function ToolFunc `json:"function"`
}

// ToolFunc 描述一个函数工具：名字、描述和 JSON Schema 参数。
type ToolFunc struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// ToolCall 是模型发起的一次完整工具调用。
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolCallFunc `json:"function"`
}

// ToolCallDelta 是流式响应中工具调用的分片，按 Index 拼接。
type ToolCallDelta struct {
	Index    int          `json:"index"`
	ID       string       `json:"id,omitempty"`
	Type     string       `json:"type,omitempty"`
	Function ToolCallFunc `json:"function"`
}

// ToolCallFunc 是函数名和 JSON 字符串形式的参数。
type ToolCallFunc struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Chatter 是"能调一次模型"的最小接口。练习里的 Agent 都依赖它，
// 测试时传 ScriptedChatter 或连到 MockServer 的 Client。
type Chatter interface {
	Complete(ctx context.Context, model string, msgs []Message, tools []Tool) (Message, error)
}

// ToolRunner 按名字执行工具，参数是 JSON 字符串。
type ToolRunner interface {
	Call(ctx context.Context, name, argsJSON string) (string, error)
}

// NewFunctionTool 构造一个 type=function 的工具定义。
func NewFunctionTool(name, desc, schema string) Tool {
	return Tool{Type: "function", Function: ToolFunc{Name: name, Description: desc, Parameters: json.RawMessage(schema)}}
}

// AssistantText 构造一条纯文本的 assistant 消息。
func AssistantText(s string) Message { return Message{Role: "assistant", Content: s} }

// AssistantCalls 构造一条带工具调用的 assistant 消息。
func AssistantCalls(calls ...ToolCall) Message { return Message{Role: "assistant", ToolCalls: calls} }

// Call 构造一个 ToolCall。
func Call(id, name, args string) ToolCall {
	return ToolCall{ID: id, Type: "function", Function: ToolCallFunc{Name: name, Arguments: args}}
}
