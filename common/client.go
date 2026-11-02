package common

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// StatusError 表示上游返回了非 200 状态码。
type StatusError struct{ Code int }

func (e *StatusError) Error() string { return fmt.Sprintf("upstream status %d", e.Code) }

// Client 是最小的流式 OpenAI 兼容客户端（就是你在 你的 Agent 项目 里写过的 Stream 的精简版）。
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient 创建客户端。
func NewClient(baseURL string) *Client { return &Client{BaseURL: baseURL, HTTP: http.DefaultClient} }

type streamFrame struct {
	Choices []struct {
		Delta struct {
			Content   string          `json:"content"`
			ToolCalls []ToolCallDelta `json:"tool_calls"`
		} `json:"delta"`
	} `json:"choices"`
}

// Complete 发一次流式请求，把分片拼成完整的 assistant 消息。实现 Chatter。
func (c *Client) Complete(ctx context.Context, model string, msgs []Message, tools []Tool) (Message, error) {
	req := ChatRequest{Model: model, Stream: true, Messages: msgs, Tools: tools}
	if len(tools) > 0 {
		req.ToolChoice = "auto"
	}
	body, err := json.Marshal(req)
	if err != nil {
		return Message{}, err
	}
	hr, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.BaseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Message{}, err
	}
	hr.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(hr)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Message{}, &StatusError{Code: resp.StatusCode}
	}
	var content strings.Builder
	slots := map[int]*ToolCall{}
	sc := bufio.NewScanner(resp.Body)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		var f streamFrame
		if err := json.Unmarshal([]byte(data), &f); err != nil {
			return Message{}, fmt.Errorf("bad frame: %w", err)
		}
		for _, ch := range f.Choices {
			content.WriteString(ch.Delta.Content)
			for _, d := range ch.Delta.ToolCalls {
				mergeDelta(slots, d)
			}
		}
	}
	if err := sc.Err(); err != nil {
		return Message{}, err
	}
	return Message{Role: "assistant", Content: content.String(), ToolCalls: assemble(slots)}, nil
}

func mergeDelta(slots map[int]*ToolCall, d ToolCallDelta) {
	tc, ok := slots[d.Index]
	if !ok {
		tc = &ToolCall{}
		slots[d.Index] = tc
	}
	if d.ID != "" {
		tc.ID = d.ID
	}
	if d.Type != "" {
		tc.Type = d.Type
	}
	if d.Function.Name != "" {
		tc.Function.Name = d.Function.Name
	}
	tc.Function.Arguments += d.Function.Arguments
}

func assemble(slots map[int]*ToolCall) []ToolCall {
	if len(slots) == 0 {
		return nil
	}
	idx := make([]int, 0, len(slots))
	for i := range slots {
		idx = append(idx, i)
	}
	sort.Ints(idx)
	out := make([]ToolCall, 0, len(idx))
	for _, i := range idx {
		out = append(out, *slots[i])
	}
	return out
}
