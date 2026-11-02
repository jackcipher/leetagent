//go:build !solution

package ex02_rate_limiter

import (
	"net/http"
	"time"
)

// Clock 抽象时间。
type Clock interface {
	Now() time.Time
}

// Bucket 是令牌桶：容量 Capacity，每秒补充 Rate 个令牌，初始是满的。
type Bucket struct {
	// TODO: 字段自行设计
}

func NewBucket(capacity, rate float64, clock Clock) *Bucket {
	// TODO
	return &Bucket{}
}

// Allow 尝试取 1 个令牌：先按经过的时间补充（不超过容量），够 1 个就扣掉返回 true。
func (b *Bucket) Allow() bool {
	// TODO
	return false
}

// Limiter 按 key（用户 ID / IP）各自一个令牌桶，并发安全。
type Limiter struct {
	Capacity float64
	Rate     float64
	Clock    Clock
	// TODO
}

func (l *Limiter) Allow(key string) bool {
	// TODO
	return false
}

// DailyQuota 是每日 token 额度：按 Clock 的 UTC 日期计，换日自动清零。并发安全。
type DailyQuota struct {
	Limit int
	Clock Clock
	// TODO
}

// Use 尝试消耗 n 个额度：used+n <= Limit 时记账并返回 true，否则不记账返回 false。
func (q *DailyQuota) Use(key string, n int) bool {
	// TODO
	return false
}

// Remaining 返回 key 今天剩余额度。
func (q *DailyQuota) Remaining(key string) int {
	// TODO
	return 0
}

// Middleware 是 HTTP 限流中间件：key 取请求头 "X-User"（为空时用 "anon"），
// 被限流时返回 429，响应头 "Retry-After: 1"，正文 "rate limited\n"；否则交给 Next。
type Middleware struct {
	Limiter *Limiter
	Next    http.Handler
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
}
