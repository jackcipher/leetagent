// Package solution 是 ex04_router 的参考答案。
package solution

import (
	"context"
	"strings"

	"leetagent/common"
)

type Route struct {
	Name        string
	Description string
	Keywords    []string
}

func hits(r Route, lower string) int {
	n := 0
	seen := map[string]bool{}
	for _, k := range r.Keywords {
		k = strings.ToLower(k)
		if k == "" || seen[k] {
			continue
		}
		seen[k] = true
		if strings.Contains(lower, k) {
			n++
		}
	}
	return n
}

func KeywordRoute(routes []Route, text string) (string, bool) {
	lower := strings.ToLower(text)
	best, bestN := "", 0
	for _, r := range routes {
		if n := hits(r, lower); n > bestN {
			best, bestN = r.Name, n
		}
	}
	return best, bestN > 0
}

func ClassifyPrompt(routes []Route) string {
	var b strings.Builder
	b.WriteString("你是一个意图分类器。根据用户输入选择最合适的类别，只输出类别名，不要输出其它内容。\n类别：\n")
	for _, r := range routes {
		b.WriteString("- " + r.Name + ": " + r.Description + "\n")
	}
	return b.String()
}

type Router struct {
	LLM     common.Chatter
	Model   string
	Routes  []Route
	Default string
}

func (r *Router) Route(ctx context.Context, text string) (string, string) {
	if name, ok := KeywordRoute(r.Routes, text); ok {
		return name, "keyword"
	}
	reply, err := r.LLM.Complete(ctx, r.Model, []common.Message{
		{Role: "system", Content: ClassifyPrompt(r.Routes)},
		{Role: "user", Content: text},
	}, nil)
	if err != nil {
		return r.Default, "default"
	}
	if name, ok := r.parse(reply.Content); ok {
		return name, "llm"
	}
	return r.Default, "default"
}

func (r *Router) parse(out string) (string, bool) {
	s := strings.ToLower(strings.Trim(strings.TrimSpace(out), "`。. \n"))
	for _, rt := range r.Routes {
		if s == rt.Name {
			return rt.Name, true
		}
	}
	found := ""
	count := 0
	for _, rt := range r.Routes {
		if strings.Contains(s, rt.Name) {
			found = rt.Name
			count++
		}
	}
	return found, count == 1
}
