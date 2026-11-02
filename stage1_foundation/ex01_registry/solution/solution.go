// Package solution 是 ex01_registry 的参考答案。
package solution

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"leetagent/common"
)

type Tool interface {
	Def() common.Tool
	Run(ctx context.Context, argsJSON string) (string, error)
}

var (
	ErrUnknownTool   = errors.New("unknown tool")
	ErrDuplicateTool = errors.New("duplicate tool")
	ErrEmptyName     = errors.New("tool name is empty")
)

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{tools: map[string]Tool{}}
}

func (r *Registry) Register(t Tool) error {
	name := t.Def().Function.Name
	if name == "" {
		return ErrEmptyName
	}
	if _, ok := r.tools[name]; ok {
		return fmt.Errorf("%w: %s", ErrDuplicateTool, name)
	}
	r.tools[name] = t
	return nil
}

func (r *Registry) Defs() []common.Tool {
	names := make([]string, 0, len(r.tools))
	for n := range r.tools {
		names = append(names, n)
	}
	sort.Strings(names)
	defs := make([]common.Tool, 0, len(names))
	for _, n := range names {
		defs = append(defs, r.tools[n].Def())
	}
	return defs
}

func (r *Registry) Call(ctx context.Context, name, argsJSON string) (string, error) {
	t, ok := r.tools[name]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownTool, name)
	}
	return t.Run(ctx, argsJSON)
}
