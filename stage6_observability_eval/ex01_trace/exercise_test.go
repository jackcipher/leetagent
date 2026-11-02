package ex01_trace

import (
	"testing"
	"time"
)

type fakeClock struct{ t time.Time }

func (f *fakeClock) Now() time.Time { return f.t }
func (f *fakeClock) Advance(ms int) { f.t = f.t.Add(time.Duration(ms) * time.Millisecond) }

func TestTraceTree(t *testing.T) {
	clk := &fakeClock{t: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	tr := &Tracer{Clock: clk}
	root := tr.Start("chat", nil)
	if root == nil {
		t.Fatal("Start returned nil")
	}
	root.SetAttr("session", "abc")
	root.SetAttr("model", "glm")
	llm := tr.Start("llm#1", root)
	llm.SetAttr("tokens", "512")
	clk.Advance(800)
	llm.End()
	tool := tr.Start("tool:get_post", root)
	clk.Advance(200)
	tool.End()
	clk.Advance(5)
	tool.End() // 重复 End 不改变
	llm2 := tr.Start("llm#2", root)
	clk.Advance(500)
	llm2.End()
	pending := tr.Start("summarize", nil)
	_ = pending
	root.End()

	want := "chat 1505ms model=glm session=abc\n" +
		"  llm#1 800ms tokens=512\n" +
		"  tool:get_post 200ms\n" +
		"  llm#2 500ms\n" +
		"summarize running\n"
	if got := tr.Render(); got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	if tool.Duration() != 200*time.Millisecond || pending.Duration() != 0 {
		t.Fatalf("durations: %v %v", tool.Duration(), pending.Duration())
	}
	if tr.Find("llm#2") != llm2 || tr.Find("nope") != nil || tr.Find("summarize") != pending {
		t.Fatal("Find broken")
	}
	if len(tr.Roots) != 2 || len(root.Children) != 3 {
		t.Fatalf("roots=%d children=%d", len(tr.Roots), len(root.Children))
	}
}
