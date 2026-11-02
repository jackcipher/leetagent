//go:build solution

package ex01_injection_guard

import "leetagent/stage7_production/ex01_injection_guard/solution"

type (
	Rule          = solution.Rule
	GuardedRunner = solution.GuardedRunner
)

var (
	ErrToolNotAllowed = solution.ErrToolNotAllowed
	DefaultRules      = solution.DefaultRules
	Detect            = solution.Detect
	WrapUntrusted     = solution.WrapUntrusted
)
