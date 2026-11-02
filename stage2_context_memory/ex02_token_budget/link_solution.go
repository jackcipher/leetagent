//go:build solution

package ex02_token_budget

import "leetagent/stage2_context_memory/ex02_token_budget/solution"

const MessageOverhead = solution.MessageOverhead

var (
	EstimateTokens = solution.EstimateTokens
	MessageTokens  = solution.MessageTokens
	TotalTokens    = solution.TotalTokens
	TrimToBudget   = solution.TrimToBudget
)
