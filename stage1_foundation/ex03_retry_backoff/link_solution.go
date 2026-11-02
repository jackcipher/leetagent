//go:build solution

package ex03_retry_backoff

import "leetagent/stage1_foundation/ex03_retry_backoff/solution"

type (
	HTTPStatusError = solution.HTTPStatusError
	Operation       = solution.Operation
	Sleeper         = solution.Sleeper
	RetryPolicy     = solution.RetryPolicy
)

var (
	IsRetryable     = solution.IsRetryable
	Do              = solution.Do
	CallWithTimeout = solution.CallWithTimeout
)
