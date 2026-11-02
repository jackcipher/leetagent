// Package solution 是 m3_research_agent 的参考答案。
package solution

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"leetagent/common"
)

var ErrMaxRounds = errors.New("max rounds exceeded")

const SystemPrompt = `你是论坛研究助手。回答问题前必须先用工具查资料：
1. 用 search_posts 搜关键词，或用 get_list 看热门帖子；
2. 用 get_post 读具体帖子正文；
3. 只根据读到的内容回答，每个事实后面用 [tid] 标注来源，例如 "pprof 可以排查泄漏 [101]"；
4. 查不到就直接说不知道，不要编造帖子编号。`

var citeRe = regexp.MustCompile(`\[(\d+)\]`)

func ExtractCitations(text string) []int {
	var out []int
	seen := map[int]bool{}
	for _, m := range citeRe.FindAllStringSubmatch(text, -1) {
		n, err := strconv.Atoi(m[1])
		if err != nil || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	return out
}

type Report struct {
	Question   string
	Answer     string
	Verified   []int
	Unverified []int
	ToolCalls  int
}

func (r Report) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 研究：%s\n\n%s\n\n## 引用\n", r.Question, r.Answer)
	if len(r.Verified)+len(r.Unverified) == 0 {
		b.WriteString("（无）\n")
	}
	for _, t := range r.Verified {
		fmt.Fprintf(&b, "- [%d] ✅\n", t)
	}
	for _, t := range r.Unverified {
		fmt.Fprintf(&b, "- [%d] ⚠️ 未验证\n", t)
	}
	return b.String()
}

type Researcher struct {
	LLM       common.Chatter
	Model     string
	Tools     common.ToolRunner
	Defs      []common.Tool
	MaxRounds int
}

func (r *Researcher) execute(ctx context.Context, c common.ToolCall) common.Message {
	out, err := r.Tools.Call(ctx, c.Function.Name, c.Function.Arguments)
	if err != nil {
		out = "error: " + err.Error()
	}
	return common.Message{Role: "tool", ToolCallID: c.ID, Content: out}
}

func (r *Researcher) verify(ctx context.Context, rep *Report) {
	for _, tid := range ExtractCitations(rep.Answer) {
		if _, err := r.Tools.Call(ctx, "get_post", fmt.Sprintf(`{"tid":%d}`, tid)); err != nil {
			rep.Unverified = append(rep.Unverified, tid)
		} else {
			rep.Verified = append(rep.Verified, tid)
		}
	}
}

func (r *Researcher) Run(ctx context.Context, question string) (Report, error) {
	rounds := r.MaxRounds
	if rounds <= 0 {
		rounds = 6
	}
	rep := Report{Question: question}
	msgs := []common.Message{{Role: "system", Content: SystemPrompt}, {Role: "user", Content: question}}
	for i := 0; i < rounds; i++ {
		reply, err := r.LLM.Complete(ctx, r.Model, msgs, r.Defs)
		if err != nil {
			return rep, err
		}
		msgs = append(msgs, reply)
		if len(reply.ToolCalls) == 0 {
			rep.Answer = reply.Content
			r.verify(ctx, &rep)
			return rep, nil
		}
		for _, c := range reply.ToolCalls {
			msgs = append(msgs, r.execute(ctx, c))
			rep.ToolCalls++
		}
	}
	return rep, ErrMaxRounds
}
