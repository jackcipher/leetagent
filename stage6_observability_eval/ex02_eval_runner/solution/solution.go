// Package solution 是 ex02_eval_runner 的参考答案。
package solution

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Case struct {
	ID             string   `json:"id"`
	Input          string   `json:"input"`
	MustContain    []string `json:"must_contain"`
	MustNotContain []string `json:"must_not_contain"`
	ExpectTool     string   `json:"expect_tool"`
}

type Output struct {
	Text      string
	ToolsUsed []string
}

type Agent interface {
	Answer(ctx context.Context, input string) (Output, error)
}

type CaseResult struct {
	ID      string
	Pass    bool
	Reasons []string
}

type Report struct {
	Total   int
	Passed  int
	Results []CaseResult
}

func LoadCases(r io.Reader) ([]Case, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	var cases []Case
	line := 0
	for sc.Scan() {
		line++
		s := strings.TrimSpace(sc.Text())
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		var c Case
		if err := json.Unmarshal([]byte(s), &c); err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		if c.ID == "" || c.Input == "" {
			return nil, fmt.Errorf("line %d: id and input are required", line)
		}
		cases = append(cases, c)
	}
	return cases, sc.Err()
}

func Grade(c Case, out Output) []string {
	var reasons []string
	lower := strings.ToLower(out.Text)
	for _, s := range c.MustContain {
		if !strings.Contains(lower, strings.ToLower(s)) {
			reasons = append(reasons, "missing: "+s)
		}
	}
	for _, s := range c.MustNotContain {
		if strings.Contains(lower, strings.ToLower(s)) {
			reasons = append(reasons, "forbidden: "+s)
		}
	}
	if c.ExpectTool != "" && !contains(out.ToolsUsed, c.ExpectTool) {
		reasons = append(reasons, "tool not used: "+c.ExpectTool)
	}
	return reasons
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func Run(ctx context.Context, agent Agent, cases []Case) Report {
	rep := Report{Total: len(cases)}
	for _, c := range cases {
		res := CaseResult{ID: c.ID}
		out, err := agent.Answer(ctx, c.Input)
		if err != nil {
			res.Reasons = []string{"error: " + err.Error()}
		} else {
			res.Reasons = Grade(c, out)
		}
		res.Pass = len(res.Reasons) == 0
		if res.Pass {
			rep.Passed++
		}
		rep.Results = append(rep.Results, res)
	}
	return rep
}

func (r Report) PassRate() float64 {
	if r.Total == 0 {
		return 0
	}
	return float64(r.Passed) / float64(r.Total)
}

func (r Report) Markdown() string {
	var b strings.Builder
	b.WriteString("| id | 结果 | 原因 |\n|---|---|---|\n")
	for _, res := range r.Results {
		if res.Pass {
			fmt.Fprintf(&b, "| %s | ✅ |  |\n", res.ID)
		} else {
			fmt.Fprintf(&b, "| %s | ❌ | %s |\n", res.ID, strings.Join(res.Reasons, "; "))
		}
	}
	fmt.Fprintf(&b, "\n通过 %d/%d（%.1f%%）\n", r.Passed, r.Total, r.PassRate()*100)
	return b.String()
}
