//go:build !solution

package ex04_cost_latency

import "time"

// Sample 是一次 LLM 调用的计时与用量。
type Sample struct {
	Start            time.Time
	FirstToken       time.Time // 收到第一个内容/工具 token 的时间
	End              time.Time
	PromptTokens     int
	CompletionTokens int
}

// Price 是每 1000 token 的价格（任意货币单位）。
type Price struct {
	PromptPer1K     float64
	CompletionPer1K float64
}

// Summary 是一批调用的统计。
type Summary struct {
	Count        int
	P50TTFT      time.Duration
	P95TTFT      time.Duration
	P50Total     time.Duration
	P95Total     time.Duration
	TokensPerSec float64 // 所有 CompletionTokens 之和 / 所有 (End-FirstToken) 之和（秒）；分母为 0 时为 0
	Cost         float64
}

// Percentile 用 nearest-rank 方法：排序后取第 ceil(p/100*n) 个（从 1 计数，至少为 1）。
// 空切片返回 0。不得修改入参。
func Percentile(values []time.Duration, p float64) time.Duration {
	// TODO
	return 0
}

// TTFT 返回首 token 延迟；Total 返回总耗时。
func TTFT(s Sample) time.Duration {
	// TODO
	return 0
}

func Total(s Sample) time.Duration {
	// TODO
	return 0
}

// Cost 计算一次调用的费用。
func (p Price) Cost(s Sample) float64 {
	// TODO
	return 0
}

// Summarize 汇总一批样本。
func Summarize(samples []Sample, price Price) Summary {
	// TODO
	return Summary{}
}
