//go:build solution

package ex04_long_term_memory

import "leetagent/stage2_context_memory/ex04_long_term_memory/solution"

const ExtractPrompt = solution.ExtractPrompt

type (
	Fact        = solution.Fact
	MemoryStore = solution.MemoryStore
)

var (
	ErrBadJSON     = solution.ErrBadJSON
	NewMemoryStore = solution.NewMemoryStore
	ParseFacts     = solution.ParseFacts
)
