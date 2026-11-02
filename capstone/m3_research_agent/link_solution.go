//go:build solution

package m3_research_agent

import "leetagent/capstone/m3_research_agent/solution"

const SystemPrompt = solution.SystemPrompt

type (
	Report     = solution.Report
	Researcher = solution.Researcher
)

var (
	ErrMaxRounds     = solution.ErrMaxRounds
	ExtractCitations = solution.ExtractCitations
)
