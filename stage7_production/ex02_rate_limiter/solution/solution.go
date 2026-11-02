// Package solution 是 ex02_rate_limiter 的参考答案。
package solution

import (
	"net/http"
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

type Bucket struct {
	capacity float64
	rate     float64
	tokens   float64
	last     time.Time
	clock    Clock
}

func NewBucket(capacity, rate float64, clock Clock) *Bucket {
	return &Bucket{capacity: capacity, rate: rate, tokens: capacity, last: clock.Now(), clock: clock}
}

func (b *Bucket) Allow() bool {
	now := b.clock.Now()
	if elapsed := now.Sub(b.last).Seconds(); elapsed > 0 {
		b.tokens += elapsed * b.rate
		if b.tokens > b.capacity {
			b.tokens = b.capacity
		}
	}
	b.last = now
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

type Limiter struct {
	Capacity float64
	Rate     float64
	Clock    Clock
	mu       sync.Mutex
	buckets  map[string]*Bucket
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.buckets == nil {
		l.buckets = map[string]*Bucket{}
	}
	b, ok := l.buckets[key]
	if !ok {
		b = NewBucket(l.Capacity, l.Rate, l.Clock)
		l.buckets[key] = b
	}
	return b.Allow()
}

type DailyQuota struct {
	Limit int
	Clock Clock
	mu    sync.Mutex
	day   string
	used  map[string]int
}

func (q *DailyQuota) roll() {
	today := q.Clock.Now().UTC().Format("2006-01-02")
	if today != q.day || q.used == nil {
		q.day = today
		q.used = map[string]int{}
	}
}

func (q *DailyQuota) Use(key string, n int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.roll()
	if q.used[key]+n > q.Limit {
		return false
	}
	q.used[key] += n
	return true
}

func (q *DailyQuota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.roll()
	return q.Limit - q.used[key]
}

type Middleware struct {
	Limiter *Limiter
	Next    http.Handler
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-User")
	if key == "" {
		key = "anon"
	}
	if !m.Limiter.Allow(key) {
		w.Header().Set("Retry-After", "1")
		http.Error(w, "rate limited", http.StatusTooManyRequests)
		return
	}
	m.Next.ServeHTTP(w, r)
}
