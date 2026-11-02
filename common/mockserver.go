package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
)

// MockReply 是 MockServer 的一条脚本回复。Status 非 0 且不是 200 时直接返回该 HTTP 状态码。
type MockReply struct {
	Status int
	Msg    Message
	// ChunkSize 控制 content 和 arguments 被切成多长的分片，默认 3 个字符（按 rune）。
	ChunkSize int
}

// ChatRequest 是 /chat/completions 的请求体。
type ChatRequest struct {
	Model      string    `json:"model"`
	Stream     bool      `json:"stream"`
	Messages   []Message `json:"messages"`
	Tools      []Tool    `json:"tools,omitempty"`
	ToolChoice string    `json:"tool_choice,omitempty"`
}

// MockServer 是离线的 OpenAI 兼容流式服务器，按脚本吐出 SSE 分片，
// 工具调用的 arguments 会被故意拆到多个 chunk 里，id/name 只在第一片出现。
type MockServer struct {
	mu       sync.Mutex
	replies  []MockReply
	Requests []ChatRequest
	srv      *httptest.Server
}

// NewMockServer 启动服务器。用完记得 Close。
func NewMockServer(replies ...MockReply) *MockServer {
	m := &MockServer{replies: replies}
	m.srv = httptest.NewServer(m)
	return m
}

// URL 返回形如 http://127.0.0.1:xxx/v1 的 BaseURL。
func (m *MockServer) URL() string { return m.srv.URL + "/v1" }

// Close 关闭服务器。
func (m *MockServer) Close() { m.srv.Close() }

// RequestCount 返回收到的请求数。
func (m *MockServer) RequestCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.Requests)
}

// ServeHTTP 实现 http.Handler。
func (m *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/chat/completions" {
		http.NotFound(w, r)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var req ChatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	m.mu.Lock()
	m.Requests = append(m.Requests, req)
	if len(m.replies) == 0 {
		m.mu.Unlock()
		http.Error(w, "script exhausted", http.StatusInternalServerError)
		return
	}
	rep := m.replies[0]
	m.replies = m.replies[1:]
	m.mu.Unlock()

	if rep.Status != 0 && rep.Status != http.StatusOK {
		http.Error(w, http.StatusText(rep.Status), rep.Status)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	for _, chunk := range BuildChunks(rep.Msg, rep.ChunkSize) {
		fmt.Fprintf(w, "data: %s\n\n", chunk)
	}
	fmt.Fprint(w, "data: [DONE]\n\n")
}

type chunkDelta struct {
	Role      string          `json:"role,omitempty"`
	Content   string          `json:"content,omitempty"`
	ToolCalls []ToolCallDelta `json:"tool_calls,omitempty"`
}

type chunkChoice struct {
	Index        int        `json:"index"`
	Delta        chunkDelta `json:"delta"`
	FinishReason *string    `json:"finish_reason"`
}

type chunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Choices []chunkChoice `json:"choices"`
}

func splitRunes(s string, n int) []string {
	if n <= 0 {
		n = 3
	}
	rs := []rune(s)
	var out []string
	for i := 0; i < len(rs); i += n {
		end := i + n
		if end > len(rs) {
			end = len(rs)
		}
		out = append(out, string(rs[i:end]))
	}
	return out
}

func encodeChunk(d chunkDelta, finish *string) string {
	b, _ := json.Marshal(chunk{ID: "mock", Object: "chat.completion.chunk", Choices: []chunkChoice{{Delta: d, FinishReason: finish}}})
	return string(b)
}

// BuildChunks 把一条完整消息拆成 SSE data 的 JSON 字符串序列（不含 "data: " 前缀）。
func BuildChunks(msg Message, size int) []string {
	out := []string{encodeChunk(chunkDelta{Role: "assistant"}, nil)}
	for _, part := range splitRunes(msg.Content, size) {
		out = append(out, encodeChunk(chunkDelta{Content: part}, nil))
	}
	for i, tc := range msg.ToolCalls {
		first := ToolCallDelta{Index: i, ID: tc.ID, Type: "function", Function: ToolCallFunc{Name: tc.Function.Name}}
		out = append(out, encodeChunk(chunkDelta{ToolCalls: []ToolCallDelta{first}}, nil))
		for _, part := range splitRunes(tc.Function.Arguments, size) {
			d := ToolCallDelta{Index: i, Function: ToolCallFunc{Arguments: part}}
			out = append(out, encodeChunk(chunkDelta{ToolCalls: []ToolCallDelta{d}}, nil))
		}
	}
	reason := "stop"
	if len(msg.ToolCalls) > 0 {
		reason = "tool_calls"
	}
	out = append(out, encodeChunk(chunkDelta{}, &reason))
	return out
}
