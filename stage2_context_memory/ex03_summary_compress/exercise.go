//go:build !solution

package ex03_summary_compress

import (
	"context"

	"leetagent/common"
)

// SummaryPrompt 是发给模型的 system 提示，要求把对话压缩成摘要。
const SummaryPrompt = "你是对话压缩器。请把下面的对话浓缩成不超过 200 字的中文摘要，保留用户的目标、已确认的事实、已调用工具得到的关键结论和未完成的事项。只输出摘要正文。"

// SummaryPrefix 是摘要消息的前缀。
const SummaryPrefix = "以下是之前对话的摘要：\n"

// RenderTranscript 把消息渲染成纯文本记录，每条一行：
//
//	user / assistant："<role>: <content>"
//	assistant 带工具调用时额外每个调用一行："assistant -> <name>(<arguments>)"（content 为空就不输出 content 行）
//	tool："tool[<tool_call_id>]: <content>"
//	system 消息跳过。
//
// 行之间用 "\n" 连接，末尾没有换行。
func RenderTranscript(msgs []common.Message) string {
	// TODO
	return ""
}

// CountTurns 返回 user 消息的个数。
func CountTurns(msgs []common.Message) int {
	// TODO
	return 0
}

// Compress 当对话轮数 > keepTurns 时，把较早的轮次压缩成摘要：
//  1. 开头的 system（如有）保留在最前；
//  2. 最近 keepTurns 轮（从倒数第 keepTurns 个 user 开始到结尾）原样保留；
//  3. 中间部分用 RenderTranscript 渲染后发给模型：
//     messages = [{system, SummaryPrompt}, {user, transcript}]，tools 传 nil；
//  4. 结果插入为一条 {role: "system", content: SummaryPrefix + 摘要(去首尾空白)}；
//
// 轮数 <= keepTurns 时原样返回且不调用模型。模型报错时返回 (nil, err)。
func Compress(ctx context.Context, c common.Chatter, model string, msgs []common.Message, keepTurns int) ([]common.Message, error) {
	// TODO
	return msgs, nil
}
