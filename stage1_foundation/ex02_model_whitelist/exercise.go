//go:build !solution

package ex02_model_whitelist

import "errors"

var (
	ErrModelNotAllowed = errors.New("model not allowed")
	ErrNoModels        = errors.New("no models configured")
	ErrBadDefault      = errors.New("default model not in allowed list")
)

// Env 抽象环境变量读取，生产传 os.Getenv，测试传假的。
type Env interface {
	Get(key string) string
}

// ModelConfig 是模型白名单配置。
type ModelConfig struct {
	Default string
	Allowed []string
}

// LoadModelConfig 从环境变量读取：
//
//	AGENT_MODELS        逗号分隔的允许列表（去掉首尾空格，忽略空项，去重保持顺序）
//	AGENT_DEFAULT_MODEL 默认模型；为空时用列表第一个
func LoadModelConfig(env Env) (ModelConfig, error) {
	// TODO
	return ModelConfig{}, errors.New("TODO")
}

// Resolve 校验前端传来的 model：空串返回默认值，在白名单里原样返回，否则 ErrModelNotAllowed。
func (c ModelConfig) Resolve(requested string) (string, error) {
	// TODO
	return "", errors.New("TODO")
}
