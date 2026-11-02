//go:build !solution

package ex02_eval_runner

import (
	"context"
	"io"
)

// Case 是一条评测用例（JSONL 的一行）。
type Case struct {
	ID             string   `json:"id"`
	Input          string   `json:"input"`
	MustContain    []string `json:"must_contain"`
	MustNotContain []string `json:"must_not_contain"`
	ExpectTool     string   `json:"expect_tool"`
}

// Output 是被测 Agent 的输出。
type Output struct {
	Text      string
	ToolsUsed []string
}

// Agent 是被测对象。
type Agent interface {
	Answer(ctx context.Context, input string) (Output, error)
}

// CaseResult 是单条用例结果。
type CaseResult struct {
	ID      string
	Pass    bool
	Reasons []string // 失败原因
}

// Report 是整体结果。
type Report struct {
	Total   int
	Passed  int
	Results []CaseResult
}

// LoadCases 读取 JSONL：跳过空行和以 "#" 开头的行；解析失败返回 error，信息以 "line N: " 开头（N 从 1 开始，按物理行号）；
// id 或 input 为空也是错误："line N: id and input are required"。
func LoadCases(r io.Reader) ([]Case, error) {
	// TODO
	return nil, nil
}

// Grade 判定一条输出，返回失败原因（空表示通过），按以下顺序：
//
//	每个 MustContain 中缺失的 -> "missing: <s>"（大小写不敏感）
//	每个 MustNotContain 中出现的 -> "forbidden: <s>"（大小写不敏感）
//	ExpectTool 非空且不在 ToolsUsed 中 -> "tool not used: <name>"
func Grade(c Case, out Output) []string {
	// TODO
	return nil
}

// Run 依次评测所有用例。Agent 出错时该用例失败，Reasons 为 ["error: " + err.Error()]。
func Run(ctx context.Context, agent Agent, cases []Case) Report {
	// TODO
	return Report{}
}

// PassRate 返回通过率（0~1），Total 为 0 时返回 0。
func (r Report) PassRate() float64 {
	// TODO
	return 0
}

// Markdown 渲染报告：
//
//	"| id | 结果 | 原因 |\n|---|---|---|\n"
//	每条 "| <id> | ✅ |  |\n" 或 "| <id> | ❌ | <原因用 "; " 连接> |\n"
//	最后 "\n通过 <Passed>/<Total>（<百分比保留一位小数>%）\n"
func (r Report) Markdown() string {
	// TODO
	return ""
}
