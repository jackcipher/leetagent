// Package solution 是 ex02_model_whitelist 的参考答案。
package solution

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrModelNotAllowed = errors.New("model not allowed")
	ErrNoModels        = errors.New("no models configured")
	ErrBadDefault      = errors.New("default model not in allowed list")
)

type Env interface {
	Get(key string) string
}

type ModelConfig struct {
	Default string
	Allowed []string
}

func LoadModelConfig(env Env) (ModelConfig, error) {
	var cfg ModelConfig
	seen := map[string]bool{}
	for _, p := range strings.Split(env.Get("AGENT_MODELS"), ",") {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		cfg.Allowed = append(cfg.Allowed, p)
	}
	if len(cfg.Allowed) == 0 {
		return ModelConfig{}, ErrNoModels
	}
	cfg.Default = strings.TrimSpace(env.Get("AGENT_DEFAULT_MODEL"))
	if cfg.Default == "" {
		cfg.Default = cfg.Allowed[0]
	}
	if !seen[cfg.Default] {
		return ModelConfig{}, fmt.Errorf("%w: %s", ErrBadDefault, cfg.Default)
	}
	return cfg, nil
}

func (c ModelConfig) Resolve(requested string) (string, error) {
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return c.Default, nil
	}
	for _, m := range c.Allowed {
		if m == requested {
			return m, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrModelNotAllowed, requested)
}
