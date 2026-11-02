package common

import (
	"context"
	"errors"
	"sync"
)

// ErrScriptExhausted 表示脚本里的回复已经用完。
var ErrScriptExhausted = errors.New("scripted chatter: no more replies")

// Reply 是脚本中的一条回复：要么返回 Msg，要么返回 Err。
type Reply struct {
	Msg Message
	Err error
}

// Call 记录 ScriptedChatter 收到的一次调用。
type CallRecord struct {
	Model string
	Msgs  []Message
	Tools []Tool
}

// ScriptedChatter 是不走网络的假模型：按顺序返回预先写好的回复，并记录每次请求。
type ScriptedChatter struct {
	mu      sync.Mutex
	Replies []Reply
	Calls   []CallRecord
}

// NewScripted 用若干消息构造 ScriptedChatter。
func NewScripted(msgs ...Message) *ScriptedChatter {
	s := &ScriptedChatter{}
	for _, m := range msgs {
		s.Replies = append(s.Replies, Reply{Msg: m})
	}
	return s
}

// Complete 实现 Chatter。
func (s *ScriptedChatter) Complete(ctx context.Context, model string, msgs []Message, tools []Tool) (Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]Message, len(msgs))
	copy(cp, msgs)
	s.Calls = append(s.Calls, CallRecord{Model: model, Msgs: cp, Tools: tools})
	if err := ctx.Err(); err != nil {
		return Message{}, err
	}
	if len(s.Replies) == 0 {
		return Message{}, ErrScriptExhausted
	}
	r := s.Replies[0]
	s.Replies = s.Replies[1:]
	return r.Msg, r.Err
}

// CallCount 返回已经被调用的次数。
func (s *ScriptedChatter) CallCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Calls)
}
