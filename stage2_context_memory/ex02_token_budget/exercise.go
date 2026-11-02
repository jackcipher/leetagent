//go:build !solution

package ex02_token_budget

import "leetagent/common"

// EstimateTokens 粗略估算 token 数（不依赖分词器）：
//
//	每个 CJK 字符（rune >= 0x2E80）算 1 个 token；
//	其它非空白字符累计，按 4 个字符 1 个 token 向上取整；空白不计。
func EstimateTokens(s string) int {
	// TODO
	return 0
}

// MessageOverhead 是每条消息的固定开销（role、分隔符等）。
const MessageOverhead = 4

// MessageTokens = MessageOverhead + EstimateTokens(Content) + 每个工具调用的 EstimateTokens(Name) + EstimateTokens(Arguments)。
func MessageTokens(m common.Message) int {
	// TODO
	return 0
}

// TotalTokens 返回所有消息的 token 之和。
func TotalTokens(msgs []common.Message) int {
	// TODO
	return 0
}

// TrimToBudget 在不超过 budget 的前提下尽量保留最近的对话：
//   - 开头的 system 消息（如果有）永远保留；
//   - 以"轮"为单位丢弃：一轮从一条 user 消息开始，到下一条 user 消息之前（包含 assistant 和 tool 消息），
//     这样不会把 tool 消息和它的 tool_calls 拆开；
//   - 从最早的轮开始丢，直到总量 <= budget；
//   - 最后一轮永远保留，即使它本身超预算。
//
// 返回新切片，不修改入参。
func TrimToBudget(msgs []common.Message, budget int) []common.Message {
	// TODO
	return msgs
}
