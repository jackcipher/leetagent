// Package solution 是 ex03_summary_compress 的参考答案。
package solution

import (
	"context"
	"fmt"
	"strings"

	"leetagent/common"
)

const SummaryPrompt = "你是对话压缩器。请把下面的对话浓缩成不超过 200 字的中文摘要，保留用户的目标、已确认的事实、已调用工具得到的关键结论和未完成的事项。只输出摘要正文。"

const SummaryPrefix = "以下是之前对话的摘要：\n"

func RenderTranscript(msgs []common.Message) string {
	var lines []string
	for _, m := range msgs {
		switch m.Role {
		case "system":
			continue
		case "tool":
			lines = append(lines, fmt.Sprintf("tool[%s]: %s", m.ToolCallID, m.Content))
		default:
			if m.Content != "" || len(m.ToolCalls) == 0 {
				lines = append(lines, m.Role+": "+m.Content)
			}
			for _, c := range m.ToolCalls {
				lines = append(lines, fmt.Sprintf("%s -> %s(%s)", m.Role, c.Function.Name, c.Function.Arguments))
			}
		}
	}
	return strings.Join(lines, "\n")
}

func CountTurns(msgs []common.Message) int {
	n := 0
	for _, m := range msgs {
		if m.Role == "user" {
			n++
		}
	}
	return n
}

// recentStart 返回倒数第 keep 个 user 消息的下标。
func recentStart(body []common.Message, keep int) int {
	seen := 0
	for i := len(body) - 1; i >= 0; i-- {
		if body[i].Role == "user" {
			seen++
			if seen == keep {
				return i
			}
		}
	}
	return 0
}

func Compress(ctx context.Context, c common.Chatter, model string, msgs []common.Message, keepTurns int) ([]common.Message, error) {
	if CountTurns(msgs) <= keepTurns {
		return msgs, nil
	}
	var head []common.Message
	body := msgs
	if len(msgs) > 0 && msgs[0].Role == "system" {
		head = msgs[:1]
		body = msgs[1:]
	}
	start := len(body)
	if keepTurns > 0 {
		start = recentStart(body, keepTurns)
	}
	old, recent := body[:start], body[start:]
	req := []common.Message{
		{Role: "system", Content: SummaryPrompt},
		{Role: "user", Content: RenderTranscript(old)},
	}
	resp, err := c.Complete(ctx, model, req, nil)
	if err != nil {
		return nil, err
	}
	out := make([]common.Message, 0, len(head)+1+len(recent))
	out = append(out, head...)
	out = append(out, common.Message{Role: "system", Content: SummaryPrefix + strings.TrimSpace(resp.Content)})
	out = append(out, recent...)
	return out, nil
}
