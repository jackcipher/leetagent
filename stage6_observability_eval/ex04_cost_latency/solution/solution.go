// Package solution 是 ex04_cost_latency 的参考答案。
package solution

import (
	"math"
	"sort"
	"time"
)

type Sample struct {
	Start            time.Time
	FirstToken       time.Time
	End              time.Time
	PromptTokens     int
	CompletionTokens int
}

type Price struct {
	PromptPer1K     float64
	CompletionPer1K float64
}

type Summary struct {
	Count        int
	P50TTFT      time.Duration
	P95TTFT      time.Duration
	P50Total     time.Duration
	P95Total     time.Duration
	TokensPerSec float64
	Cost         float64
}

type durations []time.Duration

func (d durations) Len() int           { return len(d) }
func (d durations) Less(i, j int) bool { return d[i] < d[j] }
func (d durations) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func Percentile(values []time.Duration, p float64) time.Duration {
	n := len(values)
	if n == 0 {
		return 0
	}
	sorted := append(durations(nil), values...)
	sort.Sort(sorted)
	rank := int(math.Ceil(p / 100 * float64(n)))
	if rank < 1 {
		rank = 1
	}
	if rank > n {
		rank = n
	}
	return sorted[rank-1]
}

func TTFT(s Sample) time.Duration  { return s.FirstToken.Sub(s.Start) }
func Total(s Sample) time.Duration { return s.End.Sub(s.Start) }

func (p Price) Cost(s Sample) float64 {
	return float64(s.PromptTokens)/1000*p.PromptPer1K + float64(s.CompletionTokens)/1000*p.CompletionPer1K
}

func Summarize(samples []Sample, price Price) Summary {
	sum := Summary{Count: len(samples)}
	var ttft, total []time.Duration
	var tokens int
	var gen time.Duration
	for _, s := range samples {
		ttft = append(ttft, TTFT(s))
		total = append(total, Total(s))
		tokens += s.CompletionTokens
		gen += s.End.Sub(s.FirstToken)
		sum.Cost += price.Cost(s)
	}
	sum.P50TTFT, sum.P95TTFT = Percentile(ttft, 50), Percentile(ttft, 95)
	sum.P50Total, sum.P95Total = Percentile(total, 50), Percentile(total, 95)
	if gen > 0 {
		sum.TokensPerSec = float64(tokens) / gen.Seconds()
	}
	return sum
}
