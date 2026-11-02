package ex05_hybrid_rrf

import (
	"math"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"Go语言 GMP调度!", []string{"go", "语", "言", "gmp", "调", "度"}},
		{"hello, World 2026", []string{"hello", "world", "2026"}},
		{"", nil},
	}
	for _, tc := range tests {
		if got := Tokenize(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Tokenize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func ids(rs []Ranked) []string {
	var out []string
	for _, r := range rs {
		out = append(out, r.ID)
	}
	return out
}

func TestBM25(t *testing.T) {
	docs := []Doc{
		{"d1", "Go 调度器 GMP 模型"},
		{"d2", "Rust 所有权与借用"},
		{"d3", "Go Go Go 并发"},
		{"d4", "红烧肉的做法"},
	}
	b := NewBM25(docs)
	got := b.Search("go 调度", 10)
	if !reflect.DeepEqual(ids(got), []string{"d1", "d3"}) {
		t.Fatalf("got %v", got)
	}
	if got[0].Score <= got[1].Score {
		t.Fatal("not sorted")
	}
	// 手算 d3 对 "go" 的分数：N=4 df=2 tf=3 len=5 avg=(7+7+5+6)/4
	avg := (7.0 + 7 + 5 + 6) / 4
	idf := math.Log(1 + (4-2+0.5)/(2+0.5))
	want := idf * 3 * 2.2 / (3 + 1.2*(1-0.75+0.75*5/avg))
	only := b.Search("GO go", 10)
	for _, r := range only {
		if r.ID == "d3" && math.Abs(r.Score-want) > 1e-9 {
			t.Fatalf("d3 score = %v want %v", r.Score, want)
		}
	}
	if len(b.Search("go", 1)) != 1 {
		t.Fatal("k not applied")
	}
	if len(b.Search("python", 10)) != 0 {
		t.Fatal("no match should be empty")
	}
}

func TestRRF(t *testing.T) {
	vector := []string{"a", "b", "c"}
	keyword := []string{"c", "a", "d"}
	got := RRF([][]string{vector, keyword}, 60)
	if !reflect.DeepEqual(ids(got), []string{"a", "c", "b", "d"}) {
		t.Fatalf("got %v", got)
	}
	if math.Abs(got[0].Score-(1.0/61+1.0/62)) > 1e-12 {
		t.Fatalf("score a = %v", got[0].Score)
	}
	// b 和 d 都只在一路里排第 3：同分按 ID 升序
	tie := RRF([][]string{{"x", "b"}, {"y", "a"}}, 0)
	if !reflect.DeepEqual(ids(tie), []string{"x", "y", "a", "b"}) {
		t.Fatalf("tie %v", ids(tie))
	}
}
