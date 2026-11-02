//go:build !solution

package ex01_registry

import (
	"context"
	"errors"

	"leetagent/common"
)

// Tool 是每个工具都要实现的接口。名字从 Def().Function.Name 取。
type Tool interface {
	Def() common.Tool
	Run(ctx context.Context, argsJSON string) (string, error)
}

var (
	ErrUnknownTool   = errors.New("unknown tool")
	ErrDuplicateTool = errors.New("duplicate tool")
	ErrEmptyName     = errors.New("tool name is empty")
)

// Registry 是"工具名 -> 工具"的表。
type Registry struct {
	// TODO: 你需要的字段
}

// NewRegistry 创建空注册表。
func NewRegistry() *Registry {
	// TODO
	return &Registry{}
}

// Register 注册工具。名字为空返回 ErrEmptyName，重名返回 ErrDuplicateTool。
func (r *Registry) Register(t Tool) error {
	// TODO
	return errors.New("TODO")
}

// Defs 返回所有工具定义，按名字升序。
func (r *Registry) Defs() []common.Tool {
	// TODO
	return nil
}

// Call 按名字执行工具。找不到时返回的 error 要能被 errors.Is(err, ErrUnknownTool) 识别，并带上工具名。
func (r *Registry) Call(ctx context.Context, name, argsJSON string) (string, error) {
	// TODO
	return "", errors.New("TODO")
}
