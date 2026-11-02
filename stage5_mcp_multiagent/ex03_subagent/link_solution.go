//go:build solution

package ex03_subagent

import "leetagent/stage5_mcp_multiagent/ex03_subagent/solution"

type (
	SubAgent     = solution.SubAgent
	Orchestrator = solution.Orchestrator
)

var (
	ErrMaxRounds   = solution.ErrMaxRounds
	ErrBadDelegate = solution.ErrBadDelegate
	DelegateTool   = solution.DelegateTool
)
