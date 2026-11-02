package ex04_long_term_memory

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"leetagent/common"
)

func TestUpsertAll(t *testing.T) {
	m := NewMemoryStore()
	m.Upsert(Fact{" 城市 ", " 上海 "})
	m.Upsert(Fact{"语言", "Go"})
	m.Upsert(Fact{"城市", "杭州"})
	m.Upsert(Fact{"  ", "忽略"})
	want := []Fact{{"城市", "杭州"}, {"语言", "Go"}}
	if got := m.All(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v", got)
	}
}

func TestParseFacts(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []Fact
		wantErr bool
	}{
		{"plain", `[{"key":"职业","value":"后端"}]`, []Fact{{"职业", "后端"}}, false},
		{"fence", "```json\n[{\"key\":\"a\",\"value\":\"b\"}]\n```", []Fact{{"a", "b"}}, false},
		{"empty", " [] ", []Fact{}, false},
		{"prose", "好的，我记住了", nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseFacts(tc.in)
			if tc.wantErr {
				if !errors.Is(err, ErrBadJSON) {
					t.Fatalf("err = %v", err)
				}
				return
			}
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %+v %v", got, err)
			}
		})
	}
}

func TestExtractAndRender(t *testing.T) {
	m := NewMemoryStore()
	if m.Render(0) != "" {
		t.Fatal("empty store should render empty string")
	}
	fake := common.NewScripted(common.AssistantText("```json\n[{\"key\":\"昵称\",\"value\":\"Jack\"},{\"key\":\"兴趣\",\"value\":\"Agent\"}]\n```"))
	facts, err := m.Extract(context.Background(), fake, "qwen", "我叫 Jack，最近在学 Agent")
	if err != nil || len(facts) != 2 {
		t.Fatalf("facts=%v err=%v", facts, err)
	}
	c := fake.Calls[0]
	if c.Msgs[0].Content != ExtractPrompt || c.Msgs[1].Content != "我叫 Jack，最近在学 Agent" {
		t.Fatalf("bad request %+v", c.Msgs)
	}
	if got := m.Render(0); got != "已知的用户信息：\n- 兴趣: Agent\n- 昵称: Jack" {
		t.Fatalf("Render = %q", got)
	}
	if got := m.Render(1); got != "已知的用户信息：\n- 兴趣: Agent" {
		t.Fatalf("Render(1) = %q", got)
	}
}
