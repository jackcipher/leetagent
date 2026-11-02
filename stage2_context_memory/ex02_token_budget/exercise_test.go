package ex02_token_budget

import (
	"strings"
	"testing"

	"leetagent/common"
)

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"你好世界", 4},
		{"hello", 2},
		{"hello world", 3}, // 10 个非空白字符 -> 3
		{"今天 is good", 2 + 2},
		{"   ", 0},
	}
	for _, tc := range tests {
		if got := EstimateTokens(tc.in); got != tc.want {
			t.Errorf("EstimateTokens(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMessageTokens(t *testing.T) {
	m := common.AssistantCalls(common.Call("c1", "get_list", `{"p":1}`))
	m.Content = "好的"
	// 4 + 2 + ceil(8/4)=2 + ceil(7/4)=2
	if got := MessageTokens(m); got != 10 {
		t.Fatalf("got %d want 10", got)
	}
	if TotalTokens([]common.Message{m, m}) != 20 {
		t.Fatal("TotalTokens wrong")
	}
}

func roles(ms []common.Message) string {
	var b []string
	for _, m := range ms {
		b = append(b, m.Role+":"+m.Content)
	}
	return strings.Join(b, "|")
}

func TestTrimToBudget(t *testing.T) {
	sys := common.Message{Role: "system", Content: "系统"}               // 4+2=6
	u1 := common.Message{Role: "user", Content: "一一一一一一"}              // 10
	a1 := common.AssistantCalls(common.Call("c", "x", "{}"))           // 4+1+1=6
	t1 := common.Message{Role: "tool", ToolCallID: "c", Content: "结果"} // 6
	u2 := common.Message{Role: "user", Content: "二二"}                  // 6
	a2 := common.Message{Role: "assistant", Content: "答答"}             // 6
	u3 := common.Message{Role: "user", Content: "三三三三三三三三三三三三三三三三"}    // 20
	all := []common.Message{sys, u1, a1, t1, u2, a2, u3}
	// total = 6+10+6+6+6+6+20 = 60
	tests := []struct {
		name   string
		budget int
		want   []common.Message
	}{
		{"fits", 60, all},
		{"drop first turn", 59, []common.Message{sys, u2, a2, u3}},
		{"drop two turns", 37, []common.Message{sys, u3}},
		{"last turn always kept", 1, []common.Message{sys, u3}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := append([]common.Message(nil), all...)
			got := TrimToBudget(in, tc.budget)
			if roles(got) != roles(tc.want) {
				t.Fatalf("got  %s\nwant %s", roles(got), roles(tc.want))
			}
			if roles(in) != roles(all) {
				t.Fatal("input was modified")
			}
		})
	}
	noSys := []common.Message{u1, a1, t1, u2}
	if got := TrimToBudget(noSys, 7); roles(got) != roles([]common.Message{u2}) {
		t.Fatalf("no system: %s", roles(got))
	}
}
