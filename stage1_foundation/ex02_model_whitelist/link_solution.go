//go:build solution

package ex02_model_whitelist

import "leetagent/stage1_foundation/ex02_model_whitelist/solution"

type (
	Env         = solution.Env
	ModelConfig = solution.ModelConfig
)

var (
	ErrModelNotAllowed = solution.ErrModelNotAllowed
	ErrNoModels        = solution.ErrNoModels
	ErrBadDefault      = solution.ErrBadDefault
	LoadModelConfig    = solution.LoadModelConfig
)
