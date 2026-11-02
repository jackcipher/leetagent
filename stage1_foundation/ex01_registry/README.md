# ex01 工具 Registry（Easy）

## 题目
你的 Agent 项目 现在每加一个工具要改三处：`llm.Def()`、`tools.Run` 的 switch、工具函数本身。
请实现一个 `Registry`，让"加工具"变成"写一个 struct + Register 一行"。

## 接口
```go
type Tool interface {
    Def() common.Tool
    Run(ctx context.Context, argsJSON string) (string, error)
}
func NewRegistry() *Registry
func (r *Registry) Register(t Tool) error
func (r *Registry) Defs() []common.Tool
func (r *Registry) Call(ctx context.Context, name, argsJSON string) (string, error)
```

## 要求
- 工具名取自 `t.Def().Function.Name`；为空返回 `ErrEmptyName`，重名返回 `ErrDuplicateTool`（可以用 `%w` 包装）。
- `Defs()` 按名字升序，保证每次请求里 tools 顺序稳定（便于缓存、便于测试）。
- `Call` 找不到工具时返回能被 `errors.Is(err, ErrUnknownTool)` 识别的错误，且错误信息里带工具名。

## 验收
`go test ./stage1_foundation/ex01_registry/` 全部通过。

## 搬回你自己的 Agent 项目
完成后把 Registry 放到 `server/tools/registry.go`，llm 包的 `Stream` 增加 `toolDefs []llm.Tool` 参数，删掉 `llm/tools.go`。依赖方向：controller → tools → llm，llm 不能 import tools。
