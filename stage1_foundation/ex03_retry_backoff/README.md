# ex03 超时、取消与指数退避重试（Medium）

## 题目
自建/第三方上游经常挂、偶尔 502。我们要：
1. 每次模型/工具调用都有超时（`context.WithTimeout`）。
2. 对"可能临时性"的错误按指数退避重试；对"肯定没救"的错误立即失败。

## 接口
```go
func (p RetryPolicy) Delay(attempt int) time.Duration
func IsRetryable(err error) bool
func Do(ctx context.Context, p RetryPolicy, s Sleeper, op Operation) error
func CallWithTimeout(ctx context.Context, d time.Duration, op Operation) error
```
`Operation`、`Sleeper` 都是接口——不用闭包，也方便测试里"假睡"。

## 规则
| 错误 | 是否重试 |
|---|---|
| nil | 否（成功） |
| context.Canceled / DeadlineExceeded（含被包装） | 否 |
| HTTPStatusError 429 / 5xx（含被包装） | 是 |
| HTTPStatusError 其它（400/401/404…） | 否 |
| 其它错误（连接被拒、EOF…） | 是 |

- `Delay(n) = BaseDelay * 2^(n-1)`，封顶 `MaxDelay`。
- `Do`：最多 `MaxAttempts` 次（<=0 当 1）。第 n 次失败且还能重试时先 `Sleep(ctx, Delay(n))`，Sleep 返回错误就直接返回它。最后一次失败后不再 Sleep。

## 验收
`go test ./stage1_foundation/ex03_retry_backoff/`

## 生产补充
真实系统还会加"抖动"（jitter，随机化等待避免惊群）和读取 `Retry-After` 头。流式请求**已经开始输出后**不要重试（用户会看到重复内容），这在阶段 7 的 ex03 里处理。
