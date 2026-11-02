package ex03_llm_judge

import (
	"context"
	"errors"
	"testing"

	"leetagent/common"
)

func TestParseVerdict(t *testing.T) {
	tests := []struct {
		in      string
		want    Verdict
		wantErr bool
	}{
		{`{"score":4,"reason":"基本正确"}`, Verdict{4, "基本正确"}, false},
		{"```json\n{\"score\": 5.0, \"reason\": \"完整\"}\n```", Verdict{5, "完整"}, false},
		{"评分如下：{\"score\":1,\"reason\":\"答非所问\"} 以上。", Verdict{1, "答非所问"}, false},
		{`{"score":4.5,"reason":"x"}`, Verdict{}, true},
		{`{"score":0,"reason":"x"}`, Verdict{}, true},
		{`{"score":6}`, Verdict{}, true},
		{`{"score":"4"}`, Verdict{}, true},
		{`我给 4 分`, Verdict{}, true},
	}
	for _, tc := range tests {
		got, err := ParseVerdict(tc.in)
		if tc.wantErr {
			if !errors.Is(err, ErrBadVerdict) {
				t.Errorf("%q: want ErrBadVerdict, got %v %v", tc.in, got, err)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("%q: got %+v %v", tc.in, got, err)
		}
	}
}

func TestJudgeUserPrompt(t *testing.T) {
	want := "【问题】\nQ\n\n【参考答案】\nR\n\n【待评回答】\nA"
	if got := JudgeUserPrompt("Q", "R", "A"); got != want {
		t.Fatalf("got %q", got)
	}
}

func TestJudge(t *testing.T) {
	llm := common.NewScripted(common.AssistantText("我觉得不错"), common.AssistantText(`{"score":3,"reason":"漏了要点"}`))
	v, err := Judge(context.Background(), llm, "glm", "GMP 是什么", "G、M、P 三者", "G 是 goroutine")
	if err != nil || v.Score != 3 || llm.CallCount() != 2 {
		t.Fatalf("v=%+v err=%v calls=%d", v, err, llm.CallCount())
	}
	c := llm.Calls[0]
	if c.Tools != nil || c.Msgs[0].Content != JudgeSystem || c.Msgs[1].Content != JudgeUserPrompt("GMP 是什么", "G、M、P 三者", "G 是 goroutine") {
		t.Fatalf("call = %+v", c)
	}
	bad := common.NewScripted(common.AssistantText("?"), common.AssistantText("??"), common.AssistantText(`{"score":5}`))
	if _, err := Judge(context.Background(), bad, "m", "q", "r", "a"); !errors.Is(err, ErrBadVerdict) || bad.CallCount() != 2 {
		t.Fatalf("err=%v calls=%d", err, bad.CallCount())
	}
	fail := &common.ScriptedChatter{Replies: []common.Reply{{Err: errors.New("down")}}}
	if _, err := Judge(context.Background(), fail, "m", "q", "r", "a"); err == nil || err.Error() != "down" {
		t.Fatalf("err = %v", err)
	}
}

func TestAgreement(t *testing.T) {
	tests := []struct {
		a, b []int
		tol  int
		want float64
	}{
		{[]int{5, 4, 3, 1}, []int{5, 3, 1, 1}, 0, 0.5},
		{[]int{5, 4, 3, 1}, []int{5, 3, 1, 1}, 1, 0.75},
		{[]int{1}, []int{1, 2}, 1, 0},
		{nil, nil, 0, 0},
	}
	for _, tc := range tests {
		if got := Agreement(tc.a, tc.b, tc.tol); got != tc.want {
			t.Errorf("Agreement(%v,%v,%d) = %v", tc.a, tc.b, tc.tol, got)
		}
	}
}
