//go:build solution

package ex01_session_store

import "leetagent/stage2_context_memory/ex01_session_store/solution"

type Store = solution.Store

var (
	ErrSessionNotFound = solution.ErrSessionNotFound
	Open               = solution.Open
)
