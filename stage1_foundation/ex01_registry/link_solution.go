//go:build solution

package ex01_registry

import "leetagent/stage1_foundation/ex01_registry/solution"

type (
	Tool     = solution.Tool
	Registry = solution.Registry
)

var (
	ErrUnknownTool   = solution.ErrUnknownTool
	ErrDuplicateTool = solution.ErrDuplicateTool
	ErrEmptyName     = solution.ErrEmptyName
	NewRegistry      = solution.NewRegistry
)
