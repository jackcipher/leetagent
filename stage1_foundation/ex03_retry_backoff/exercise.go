//go:build !solution

package ex03_retry_backoff

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// HTTPStatusError 表示上游返回了某个 HTTP 状态码。
type HTTPStatusError struct{ Code int }

func (e *HTTPStatusError) Error() string { return fmt.Sprintf("http status %d", e.Code) }

// Operation 是一次可重试的操作（例如一次模型调用）。
type Operation interface {
	Run(ctx context.Context) error
}

// Sleeper 抽象"等待"，测试里用假的，不真的睡。
type Sleeper interface {
	Sleep(ctx context.Context, d time.Duration) error
}

// RetryPolicy 是重试策略。
type RetryPolicy struct {
	MaxAttempts int           // 总尝试次数（含第一次），<=0 视为 1
	BaseDelay   time.Duration // 第 1 次重试前的等待
	MaxDelay    time.Duration // 等待上限
}

// Delay 返回第 attempt 次重试前的等待时间：BaseDelay * 2^(attempt-1)，不超过 MaxDelay。attempt 从 1 开始。
func (p RetryPolicy) Delay(attempt int) time.Duration {
	// TODO
	return 0
}

// IsRetryable 判断错误是否值得重试：
// context.Canceled / DeadlineExceeded -> false；HTTPStatusError 429 或 >=500 -> true，其它状态码 -> false；
// nil -> false；其它错误（网络错误等）-> true。
func IsRetryable(err error) bool {
	// TODO
	return false
}

// Do 按策略执行 op：成功立即返回 nil；不可重试的错误立即返回；
// 重试前调用 s.Sleep(ctx, p.Delay(n))，Sleep 出错（ctx 取消）就返回该错误；用完次数返回最后一次错误。
func Do(ctx context.Context, p RetryPolicy, s Sleeper, op Operation) error {
	// TODO
	return errors.New("TODO")
}

// CallWithTimeout 给 op 加上超时：用 context.WithTimeout 派生 ctx 后调用 op.Run。
func CallWithTimeout(ctx context.Context, d time.Duration, op Operation) error {
	// TODO
	return errors.New("TODO")
}
