package ex04_cost_latency

import (
	"math"
	"testing"
	"time"
)

func ms(n int) time.Duration { return time.Duration(n) * time.Millisecond }

func TestPercentile(t *testing.T) {
	vals := []time.Duration{ms(500), ms(100), ms(400), ms(200), ms(300), ms(900), ms(700), ms(600), ms(800), ms(1000)}
	tests := []struct {
		p    float64
		want time.Duration
	}{{50, ms(500)}, {95, ms(1000)}, {90, ms(900)}, {10, ms(100)}, {0, ms(100)}, {100, ms(1000)}}
	for _, tc := range tests {
		if got := Percentile(vals, tc.p); got != tc.want {
			t.Errorf("p%v = %v want %v", tc.p, got, tc.want)
		}
	}
	if vals[0] != ms(500) {
		t.Fatal("input modified")
	}
	if Percentile(nil, 50) != 0 {
		t.Fatal("empty should be 0")
	}
}

func TestSummarize(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	mk := func(ttft, total, pt, ct int) Sample {
		return Sample{Start: base, FirstToken: base.Add(ms(ttft)), End: base.Add(ms(total)), PromptTokens: pt, CompletionTokens: ct}
	}
	samples := []Sample{mk(200, 1200, 1000, 100), mk(400, 2400, 2000, 200), mk(1500, 2500, 500, 100)}
	if TTFT(samples[0]) != ms(200) || Total(samples[0]) != ms(1200) {
		t.Fatal("TTFT/Total")
	}
	price := Price{PromptPer1K: 0.5, CompletionPer1K: 2}
	if c := price.Cost(samples[0]); math.Abs(c-0.7) > 1e-9 {
		t.Fatalf("cost = %v", c)
	}
	s := Summarize(samples, price)
	if s.Count != 3 || s.P50TTFT != ms(400) || s.P95TTFT != ms(1500) || s.P50Total != ms(2400) || s.P95Total != ms(2500) {
		t.Fatalf("summary = %+v", s)
	}
	// tokens 400 / gen (1.0+2.0+1.0)s = 100
	if math.Abs(s.TokensPerSec-100) > 1e-9 {
		t.Fatalf("tps = %v", s.TokensPerSec)
	}
	if math.Abs(s.Cost-(0.7+1.4+0.45)) > 1e-9 {
		t.Fatalf("cost = %v", s.Cost)
	}
	if z := Summarize(nil, price); z.TokensPerSec != 0 || z.Count != 0 {
		t.Fatalf("empty = %+v", z)
	}
}
