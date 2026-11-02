//go:build solution

package ex02_rate_limiter

import "leetagent/stage7_production/ex02_rate_limiter/solution"

type (
	Clock      = solution.Clock
	Bucket     = solution.Bucket
	Limiter    = solution.Limiter
	DailyQuota = solution.DailyQuota
	Middleware = solution.Middleware
)

var NewBucket = solution.NewBucket
