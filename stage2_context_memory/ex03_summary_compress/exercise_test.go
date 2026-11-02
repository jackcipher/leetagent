package ex03_summary_compress

import (
	"context"
	"errors"
	"testing"

	"leetagent/common"
)

func convo() []common.Message {
	return []common.Message{
		{Role: "system", Content: "你是论坛助手"},
		{Role: "user", Content: "看第一页"},
		common.AssistantCalls(common.Call("c1", "get_list", `{"page":1}`)),
		{Role: "tool", ToolCallID: "c1", Content: "3 个帖子"},
		{Role: "assistant", Content: "有 3 个帖子"},
		{Role: "user", Content: "打开第二个"},
		{Role: "assistant", Content: "帖子讲的是 Go"},
		{Role: "user", Content: "总结一下"},
	}
}

func TestRenderTranscript(t *testing.T) {
	got := RenderTranscript(convo()[:6])
	want := "user: 看第一页\nassistant -> get_list({\"page\":1})\ntool[c1]: 3 个帖子\nassistant: 有 3 个帖子\nuser: 打开第二个"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	if CountTurns(convo()) != 3 {
		t.Fatal("CountTurns wrong")
	}
}

func TestCompress(t *testing.T) {
	fake := common.NewScripted(common.AssistantText("  用户在看论坛第一页，共 3 帖。  "))
	got, err := Compress(context.Background(), fake, "glm", convo(), 2)
	if err != nil {
		t.Fatal(err)
	}
	if fake.CallCount() != 1 {
		t.Fatalf("model calls = %d", fake.CallCount())
	}
	call := fake.Calls[0]
	if call.Model != "glm" || len(call.Msgs) != 2 || call.Msgs[0].Content != SummaryPrompt || call.Tools != nil {
		t.Fatalf("bad summary request: %+v", call)
	}
	wantTranscript := "user: 看第一页\nassistant -> get_list({\"page\":1})\ntool[c1]: 3 个帖子\nassistant: 有 3 个帖子"
	if call.Msgs[1].Content != wantTranscript {
		t.Fatalf("transcript:\n%s", call.Msgs[1].Content)
	}
	if len(got) != 5 {
		t.Fatalf("len = %d: %+v", len(got), got)
	}
	if got[0].Content != "你是论坛助手" || got[1].Role != "system" || got[1].Content != SummaryPrefix+"用户在看论坛第一页，共 3 帖。" ||
		got[2].Content != "打开第二个" || got[4].Content != "总结一下" {
		t.Fatalf("bad result: %+v", got)
	}
}

func TestCompressNoop(t *testing.T) {
	fake := common.NewScripted()
	in := convo()
	got, err := Compress(context.Background(), fake, "m", in, 3)
	if err != nil || len(got) != len(in) || fake.CallCount() != 0 {
		t.Fatalf("should be noop: len=%d err=%v calls=%d", len(got), err, fake.CallCount())
	}
}

func TestCompressError(t *testing.T) {
	fake := &common.ScriptedChatter{Replies: []common.Reply{{Err: errors.New("EOF")}}}
	if _, err := Compress(context.Background(), fake, "m", convo(), 1); err == nil {
		t.Fatal("want error")
	}
}
