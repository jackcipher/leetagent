package ex01_session_store

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"leetagent/common"
)

func TestStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "chat.db")
	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	ctx := context.Background()
	id, err := s.CreateSession(ctx, "论坛调研")
	if err != nil || id <= 0 {
		t.Fatalf("CreateSession = %d, %v", id, err)
	}
	other, _ := s.CreateSession(ctx, "另一个")
	msgs := []common.Message{
		{Role: "user", Content: "看看第一页"},
		common.AssistantCalls(common.Call("call_1", "get_list", `{"page":1}`)),
		{Role: "tool", ToolCallID: "call_1", Content: `[{"tid":"1"}]`},
		{Role: "assistant", Content: "第一页有 1 个帖子"},
	}
	for _, m := range msgs {
		if err := s.Append(ctx, id, m); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	_ = s.Append(ctx, other, common.Message{Role: "user", Content: "别混进来"})
	got, err := s.History(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, msgs) {
		t.Fatalf("History mismatch\n got %+v\nwant %+v", got, msgs)
	}
	s.Close()

	// 重新打开，数据还在
	s2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()
	got, _ = s2.History(ctx, id)
	if len(got) != 4 {
		t.Fatalf("after reopen got %d messages", len(got))
	}
}

func TestNotFound(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "x.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Append(context.Background(), 999, common.Message{Role: "user"}); !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("Append err = %v", err)
	}
	if _, err := s.History(context.Background(), 999); !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("History err = %v", err)
	}
}
