package ex04_approval_gate

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"leetagent/common"
)

type echoRunner struct{ ran []string }

func (e *echoRunner) Call(ctx context.Context, name, args string) (string, error) {
	e.ran = append(e.ran, name)
	return "done:" + name, nil
}

type fixedApprover struct {
	answers map[string]bool
	err     error
	seen    []common.ToolCall
}

func (f *fixedApprover) Approve(ctx context.Context, c common.ToolCall) (bool, error) {
	f.seen = append(f.seen, c)
	return f.answers[c.Function.Name], f.err
}

func TestGate(t *testing.T) {
	inner := &echoRunner{}
	ap := &fixedApprover{answers: map[string]bool{"post_reply": true, "delete_post": false}}
	g := &Gate{Inner: inner, Approver: ap, Dangerous: map[string]bool{"post_reply": true, "delete_post": true}}
	var _ common.ToolRunner = g

	tests := []struct {
		name    string
		want    string
		wantErr error
	}{
		{"get_list", "done:get_list", nil},
		{"post_reply", "done:post_reply", nil},
		{"delete_post", "", ErrDenied},
	}
	for _, tc := range tests {
		got, err := g.Call(context.Background(), tc.name, `{"id":1}`)
		if got != tc.want || !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: got %q %v", tc.name, got, err)
		}
	}
	if !reflect.DeepEqual(inner.ran, []string{"get_list", "post_reply"}) {
		t.Fatalf("ran = %v", inner.ran)
	}
	if !reflect.DeepEqual(g.Log, []string{"get_list: auto", "post_reply: approved", "delete_post: denied"}) {
		t.Fatalf("log = %v", g.Log)
	}
	if ap.seen[0].Function.Arguments != `{"id":1}` || ap.seen[0].Function.Name != "post_reply" {
		t.Fatalf("approver saw %+v", ap.seen[0])
	}

	ap.err = errors.New("ui gone")
	if _, err := g.Call(context.Background(), "post_reply", "{}"); err == nil || err.Error() != "ui gone" || len(inner.ran) != 2 {
		t.Fatalf("approver error must block: %v", err)
	}
}

func TestChanApprover(t *testing.T) {
	a := NewChanApprover()
	type outcome struct {
		ok  bool
		err error
	}
	done := make(chan outcome, 2)
	go func() {
		ok, err := a.Approve(context.Background(), common.Call("", "post_reply", "{}"))
		done <- outcome{ok, err}
	}()
	p := next(t, a)
	if p.ID != "req-1" || p.Call.Function.Name != "post_reply" {
		t.Fatalf("pending = %+v", p)
	}
	if a.Resolve("req-99", true) {
		t.Fatal("unknown id should return false")
	}
	if !a.Resolve(p.ID, true) {
		t.Fatal("resolve failed")
	}
	if o := <-done; !o.ok || o.err != nil {
		t.Fatalf("outcome = %+v", o)
	}
	if a.Resolve(p.ID, true) {
		t.Fatal("double resolve should return false")
	}

	go func() {
		ok, err := a.Approve(context.Background(), common.Call("", "delete_post", "{}"))
		done <- outcome{ok, err}
	}()
	p = next(t, a)
	if p.ID != "req-2" {
		t.Fatalf("id = %s", p.ID)
	}
	a.Resolve(p.ID, false)
	if o := <-done; o.ok || o.err != nil {
		t.Fatalf("outcome = %+v", o)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if a.Requests() == nil {
		t.Fatal("Requests() is nil")
	}
	go func() { <-a.Requests() }()
	ok, err := a.Approve(ctx, common.Call("", "x", "{}"))
	if ok || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout: %v %v", ok, err)
	}
}

func next(t *testing.T, a *ChanApprover) Pending {
	t.Helper()
	select {
	case p := <-a.Requests():
		return p
	case <-time.After(time.Second):
		t.Fatal("no pending request within 1s")
	}
	return Pending{}
}
