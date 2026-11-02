//go:build !solution

// Package m3_research_agent 是毕业项目第 3 步：会查论坛、会引用、引用可验证的研究 Agent。
package m3_research_agent

import (
	"context"
	"errors"

	"leetagent/common"
)

// ErrMaxRounds 表示超过最大轮数仍没有最终回答。
var ErrMaxRounds = errors.New("max rounds exceeded")

// SystemPrompt 是研究 Agent 的系统提示词（测试会检查它包含 "[tid]" 这个引用格式说明）。
const SystemPrompt = "TODO"

// ExtractCitations 找出文本里所有 "[数字]" 形式的引用，按首次出现顺序去重返回。
// "[101]" ✓  "[ 101 ]" ✗  "[abc]" ✗  "[]" ✗
func ExtractCitations(text string) []int {
	// TODO
	return nil
}

// Report 是研究结果。
type Report struct {
	Question   string
	Answer     string
	Verified   []int // 用 get_post 确认存在的引用
	Unverified []int // get_post 报错的引用（模型编造的）
	ToolCalls  int   // 研究过程中执行的工具调用次数（不含验证）
}

// Markdown 渲染报告，格式严格如下（每行以 \n 结尾）：
//
//	# 研究：<Question>
//
//	<Answer>
//
//	## 引用
//	- [101] ✅
//	- [999] ⚠️ 未验证
//
// 引用部分先列 Verified 再列 Unverified；没有任何引用时 "## 引用" 下面写一行 "（无）"。
func (r Report) Markdown() string {
	// TODO
	return ""
}

// Researcher 组合 LLM 和工具。
type Researcher struct {
	LLM       common.Chatter
	Model     string
	Tools     common.ToolRunner
	Defs      []common.Tool
	MaxRounds int // <=0 视为 6
}

// Run：
//  1. 消息 = [system: SystemPrompt, user: question]；
//  2. 循环调用 LLM（带 Defs）：有 tool_calls 就逐个执行，结果以 role=tool、ToolCallID 对应追加；
//     工具出错时内容写 "error: <err>"（让模型自己纠正），不要中断；
//     没有 tool_calls 时，Content 即最终回答；
//  3. 超过 MaxRounds 次 LLM 调用仍没回答 → ErrMaxRounds；LLM 本身出错直接返回；
//  4. 对回答里的每个引用调 Tools.Call(ctx, "get_post", `{"tid":N}`) 验证。
func (r *Researcher) Run(ctx context.Context, question string) (Report, error) {
	// TODO
	return Report{}, nil
}
