//go:build !solution

package ex04_long_term_memory

import (
	"context"
	"errors"

	"leetagent/common"
)

// ExtractPrompt 让模型从用户消息里抽取长期事实，输出 JSON 数组。
const ExtractPrompt = `从用户的话里抽取值得长期记住的事实（偏好、身份、长期目标）。只输出 JSON 数组，形如 [{"key":"城市","value":"上海"}]；没有就输出 []。`

var ErrBadJSON = errors.New("model output is not a JSON fact array")

// Fact 是一条长期记忆。
type Fact struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MemoryStore 是按 key 覆盖写入的记忆表（内存版）。
type MemoryStore struct {
	// TODO
}

// NewMemoryStore 创建空表。
func NewMemoryStore() *MemoryStore {
	// TODO
	return &MemoryStore{}
}

// Upsert 写入事实：key/value 去首尾空白，key 为空忽略，同 key 覆盖。
func (m *MemoryStore) Upsert(f Fact) {
	// TODO
}

// All 返回全部事实，按 key 升序。
func (m *MemoryStore) All() []Fact {
	// TODO
	return nil
}

// ParseFacts 解析模型输出：容忍 ```json ... ``` 代码块包裹和前后空白；不是合法 JSON 数组返回 ErrBadJSON。
func ParseFacts(output string) ([]Fact, error) {
	// TODO
	return nil, nil
}

// Extract 调用模型抽取事实并写入 store，返回写入的事实（解析结果）。
// 请求：[{system, ExtractPrompt}, {user, userMsg}]，tools 为 nil。
func (m *MemoryStore) Extract(ctx context.Context, c common.Chatter, model, userMsg string) ([]Fact, error) {
	// TODO
	return nil, nil
}

// Render 生成注入 system prompt 的文本："已知的用户信息：\n- k: v\n- k: v"（按 key 升序，最多 max 条，max<=0 不限）。
// 没有任何事实时返回空串。
func (m *MemoryStore) Render(max int) string {
	// TODO
	return ""
}
