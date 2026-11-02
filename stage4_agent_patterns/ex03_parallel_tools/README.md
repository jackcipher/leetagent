# ex03 并行工具调用（Medium）

## 题目
模型一次可能返回多个 `tool_calls`（例如"同时查 3 个帖子"）。串行执行慢；无限并发会打爆论坛 API。
实现带并发上限、保持顺序、支持取消的 `ExecuteAll`。

## 接口
```go
type Result struct { CallID, Content string; Err error }
func ExecuteAll(ctx, runner common.ToolRunner, calls []common.ToolCall, maxConcurrency int) []Result
func ToMessages(results []Result) []common.Message
```

## 要求
- 并发峰值 ≤ maxConcurrency（≤0 当作 1）
- 结果顺序 = calls 顺序（OpenAI 协议要求每个 tool_call_id 都有对应 tool 消息）
- 取消后未开始的调用不执行，Err = ctx.Err()
- **不用闭包**：`go job{...}.run()` 这种具名方法写法

## 验收
`go test -race ./stage4_agent_patterns/ex03_parallel_tools/`（加 -race 检查数据竞争；需要 cgo，若不可用去掉 -race）
