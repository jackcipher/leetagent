//go:build solution

package ex03_llm_judge

import "leetagent/stage6_observability_eval/ex03_llm_judge/solution"

const JudgeSystem = solution.JudgeSystem

type Verdict = solution.Verdict

var (
	ErrBadVerdict   = solution.ErrBadVerdict
	ParseVerdict    = solution.ParseVerdict
	JudgeUserPrompt = solution.JudgeUserPrompt
	Judge           = solution.Judge
	Agreement       = solution.Agreement
)
