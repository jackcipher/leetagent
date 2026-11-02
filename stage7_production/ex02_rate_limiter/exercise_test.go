package ex02_rate_limiter

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (f *fakeClock) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.t
}

func (f *fakeClock) Advance(d time.Duration) {
	f.mu.Lock()
	f.t = f.t.Add(d)
	f.mu.Unlock()
}

func newClock() *fakeClock { return &fakeClock{t: time.Date(2026, 3, 1, 23, 59, 0, 0, time.UTC)} }

func TestBucket(t *testing.T) {
	clk := newClock()
	b := NewBucket(3, 2, clk) // 3 个令牌，每秒补 2 个
	steps := []struct {
		advance time.Duration
		want    bool
	}{
		{0, true}, {0, true}, {0, true}, {0, false},
		{250 * time.Millisecond, false}, // 0.5 个
		{250 * time.Millisecond, true},  // 1.0 个
		{0, false},
		{10 * time.Second, true}, // 补满到 3，不会超
		{0, true}, {0, true}, {0, false},
	}
	for i, s := range steps {
		clk.Advance(s.advance)
		if got := b.Allow(); got != s.want {
			t.Fatalf("step %d: got %v want %v", i, got, s.want)
		}
	}
}

func TestLimiterPerKey(t *testing.T) {
	clk := newClock()
	l := &Limiter{Capacity: 1, Rate: 1, Clock: clk}
	if !l.Allow("alice") || l.Allow("alice") {
		t.Fatal("alice: 1 then deny")
	}
	if !l.Allow("bob") {
		t.Fatal("bob has his own bucket")
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	allowed := 0
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow("carol") {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if allowed != 1 {
		t.Fatalf("concurrent allowed = %d", allowed)
	}
}

func TestDailyQuota(t *testing.T) {
	clk := newClock()
	q := &DailyQuota{Limit: 1000, Clock: clk}
	if !q.Use("u", 600) || q.Use("u", 500) || !q.Use("u", 400) {
		t.Fatal("quota accounting")
	}
	if q.Remaining("u") != 0 || q.Remaining("v") != 1000 {
		t.Fatalf("remaining = %d %d", q.Remaining("u"), q.Remaining("v"))
	}
	clk.Advance(2 * time.Minute) // 跨过 UTC 午夜
	if q.Remaining("u") != 1000 || !q.Use("u", 1000) {
		t.Fatal("quota should reset on new day")
	}
}

type okHandler struct{}

func (okHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }

func TestMiddleware(t *testing.T) {
	m := &Middleware{Limiter: &Limiter{Capacity: 1, Rate: 0.001, Clock: newClock()}, Next: okHandler{}}
	do := func(user string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/api/chat", nil)
		if user != "" {
			req.Header.Set("X-User", user)
		}
		rec := httptest.NewRecorder()
		m.ServeHTTP(rec, req)
		return rec
	}
	if r := do("jack"); r.Code != 200 || r.Body.String() != "ok" {
		t.Fatalf("first = %d %q", r.Code, r.Body.String())
	}
	r := do("jack")
	if r.Code != 429 || r.Header().Get("Retry-After") != "1" || r.Body.String() != "rate limited\n" {
		t.Fatalf("second = %d %q %v", r.Code, r.Body.String(), r.Header())
	}
	if do("").Code != 200 || do("").Code != 429 {
		t.Fatal("anon shares one bucket")
	}
}
