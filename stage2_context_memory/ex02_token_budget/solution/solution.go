// Package solution 是 ex02_token_budget 的参考答案。
package solution

import (
	"unicode"

	"leetagent/common"
)

const MessageOverhead = 4

func EstimateTokens(s string) int {
	cjk, other := 0, 0
	for _, r := range s {
		switch {
		case unicode.IsSpace(r):
		case r >= 0x2E80:
			cjk++
		default:
			other++
		}
	}
	return cjk + (other+3)/4
}

func MessageTokens(m common.Message) int {
	n := MessageOverhead + EstimateTokens(m.Content)
	for _, c := range m.ToolCalls {
		n += EstimateTokens(c.Function.Name) + EstimateTokens(c.Function.Arguments)
	}
	return n
}

func TotalTokens(msgs []common.Message) int {
	n := 0
	for _, m := range msgs {
		n += MessageTokens(m)
	}
	return n
}

// splitTurns 返回每一轮的起始下标（相对于 body）。
func splitTurns(body []common.Message) []int {
	var starts []int
	for i, m := range body {
		if m.Role == "user" || i == 0 {
			starts = append(starts, i)
		}
	}
	return starts
}

func TrimToBudget(msgs []common.Message, budget int) []common.Message {
	var head []common.Message
	body := msgs
	if len(msgs) > 0 && msgs[0].Role == "system" {
		head = msgs[:1]
		body = msgs[1:]
	}
	starts := splitTurns(body)
	cut := 0
	for k := 0; k < len(starts)-1; k++ {
		if TotalTokens(head)+TotalTokens(body[starts[k]:]) <= budget {
			break
		}
		cut = starts[k+1]
	}
	out := make([]common.Message, 0, len(head)+len(body)-cut)
	out = append(out, head...)
	out = append(out, body[cut:]...)
	return out
}
