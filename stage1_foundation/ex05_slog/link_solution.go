//go:build solution

package ex05_slog

import "leetagent/stage1_foundation/ex05_slog/solution"

type RoundStats = solution.RoundStats

var (
	NewJSONLogger = solution.NewJSONLogger
	LogRound      = solution.LogRound
	WithRequest   = solution.WithRequest
)
