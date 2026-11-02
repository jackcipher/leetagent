// Package solution 是 ex04_approval_gate 的参考答案。
package solution

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"leetagent/common"
)

var ErrDenied = errors.New("tool call denied by human")

type Approver interface {
	Approve(ctx context.Context, call common.ToolCall) (bool, error)
}

type Gate struct {
	Inner     common.ToolRunner
	Approver  Approver
	Dangerous map[string]bool
	Log       []string
}

func (g *Gate) Call(ctx context.Context, name, args string) (string, error) {
	if !g.Dangerous[name] {
		g.Log = append(g.Log, name+": auto")
		return g.Inner.Call(ctx, name, args)
	}
	ok, err := g.Approver.Approve(ctx, common.Call("", name, args))
	if err != nil {
		return "", err
	}
	if !ok {
		g.Log = append(g.Log, name+": denied")
		return "", ErrDenied
	}
	g.Log = append(g.Log, name+": approved")
	return g.Inner.Call(ctx, name, args)
}

type Pending struct {
	ID   string
	Call common.ToolCall
}

type ChanApprover struct {
	mu      sync.Mutex
	seq     int
	reqs    chan Pending
	waiting map[string]chan bool
}

func NewChanApprover() *ChanApprover {
	return &ChanApprover{reqs: make(chan Pending, 16), waiting: map[string]chan bool{}}
}

func (a *ChanApprover) Requests() <-chan Pending { return a.reqs }

func (a *ChanApprover) Approve(ctx context.Context, call common.ToolCall) (bool, error) {
	a.mu.Lock()
	a.seq++
	id := fmt.Sprintf("req-%d", a.seq)
	answer := make(chan bool, 1)
	a.waiting[id] = answer
	a.mu.Unlock()
	defer a.forget(id)

	select {
	case a.reqs <- Pending{ID: id, Call: call}:
	case <-ctx.Done():
		return false, ctx.Err()
	}
	select {
	case ok := <-answer:
		return ok, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (a *ChanApprover) forget(id string) {
	a.mu.Lock()
	delete(a.waiting, id)
	a.mu.Unlock()
}

func (a *ChanApprover) Resolve(id string, ok bool) bool {
	a.mu.Lock()
	ch, found := a.waiting[id]
	if found {
		delete(a.waiting, id)
	}
	a.mu.Unlock()
	if !found {
		return false
	}
	ch <- ok
	return true
}
