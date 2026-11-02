//go:build solution

package ex02_eval_runner

import "leetagent/stage6_observability_eval/ex02_eval_runner/solution"

type (
	Case       = solution.Case
	Output     = solution.Output
	Agent      = solution.Agent
	CaseResult = solution.CaseResult
	Report     = solution.Report
)

var (
	LoadCases = solution.LoadCases
	Grade     = solution.Grade
	Run       = solution.Run
)
