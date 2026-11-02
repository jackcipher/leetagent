// Package solution 是 ex04_long_term_memory 的参考答案。
package solution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"leetagent/common"
)

const ExtractPrompt = `从用户的话里抽取值得长期记住的事实（偏好、身份、长期目标）。只输出 JSON 数组，形如 [{"key":"城市","value":"上海"}]；没有就输出 []。`

var ErrBadJSON = errors.New("model output is not a JSON fact array")

type Fact struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MemoryStore struct {
	facts map[string]string
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{facts: map[string]string{}} }

func (m *MemoryStore) Upsert(f Fact) {
	k := strings.TrimSpace(f.Key)
	if k == "" {
		return
	}
	m.facts[k] = strings.TrimSpace(f.Value)
}

func (m *MemoryStore) All() []Fact {
	keys := make([]string, 0, len(m.facts))
	for k := range m.facts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]Fact, 0, len(keys))
	for _, k := range keys {
		out = append(out, Fact{Key: k, Value: m.facts[k]})
	}
	return out
}

// StripFence 去掉 ``` 或 ```json 代码块包裹。
func StripFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```")
	if i := strings.Index(s, "\n"); i >= 0 {
		s = s[i+1:]
	}
	s = strings.TrimSuffix(strings.TrimSpace(s), "```")
	return strings.TrimSpace(s)
}

func ParseFacts(output string) ([]Fact, error) {
	var facts []Fact
	if err := json.Unmarshal([]byte(StripFence(output)), &facts); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadJSON, err)
	}
	if facts == nil {
		facts = []Fact{}
	}
	return facts, nil
}

func (m *MemoryStore) Extract(ctx context.Context, c common.Chatter, model, userMsg string) ([]Fact, error) {
	resp, err := c.Complete(ctx, model, []common.Message{
		{Role: "system", Content: ExtractPrompt},
		{Role: "user", Content: userMsg},
	}, nil)
	if err != nil {
		return nil, err
	}
	facts, err := ParseFacts(resp.Content)
	if err != nil {
		return nil, err
	}
	for _, f := range facts {
		m.Upsert(f)
	}
	return facts, nil
}

func (m *MemoryStore) Render(max int) string {
	all := m.All()
	if len(all) == 0 {
		return ""
	}
	if max > 0 && len(all) > max {
		all = all[:max]
	}
	var b strings.Builder
	b.WriteString("已知的用户信息：")
	for _, f := range all {
		b.WriteString("\n- " + f.Key + ": " + f.Value)
	}
	return b.String()
}
