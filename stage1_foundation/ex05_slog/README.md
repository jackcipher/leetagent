# ex05 结构化日志 slog（Easy）

## 题目
`fmt.Println` 打日志没法检索。生产里日志要是 JSON，一行一个事件，字段固定，才能用 jq / Loki / ELK 查询：
"昨天哪些请求调了 fetch_one 超过 3 秒？"

## 接口
```go
func NewJSONLogger(w io.Writer, level slog.Level) *slog.Logger
func LogRound(l *slog.Logger, s RoundStats)
func WithRequest(l *slog.Logger, requestID, user string) *slog.Logger
```

## 要求
- 消息 `agent round`；属性：`request_id round model duration_ms tool_count tools`（tools 用逗号拼）。
- 有 Err 时级别 ERROR，并带 `err` 字符串；否则 INFO。
- `WithRequest` 用 `l.With(...)` 派生子 logger，后续每条日志自动带 request_id 和 user。

## 验收
`go test ./stage1_foundation/ex05_slog/`

## 用法示例
```bash
go run ./server | jq 'select(.msg=="agent round" and .duration_ms > 3000)'
```
