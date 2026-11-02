package m2_research_index

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
		{"Go 语言", []string{"go", "语言"}},
		{"显存估算", []string{"显存", "存估", "估算"}},
		{"7B模型，int4量化!", []string{"7b", "模型", "int4", "量化"}},
		{"用 pprof", []string{"用", "pprof"}},
		{"  ，。", nil},
	}
	for _, tc := range tests {
		got := Tokenize(tc.in)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Tokenize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestEmbed(t *testing.T) {
	e := HashEmbedder{Dim: 64}
	v := e.Embed("goroutine 泄漏 goroutine")
	if len(v) != 64 {
		t.Fatalf("len = %d", len(v))
	}
	var n float64
	for _, x := range v {
		n += x * x
	}
	if math.Abs(n-1) > 1e-9 {
		t.Fatalf("norm^2 = %v, want 1", n)
	}
	if !reflect.DeepEqual(v, e.Embed("goroutine 泄漏 goroutine")) {
		t.Fatal("embedding must be deterministic")
	}
	z := e.Embed("，。")
	if len(z) != 64 || z[0] != 0 {
		t.Fatal("empty text → zero vector of Dim")
	}
}

func TestSearch(t *testing.T) {
	ix := NewIndex(256)
	if ix == nil {
		t.Fatal("NewIndex returned nil")
	}
	ix.Add(101, "goroutine 泄漏 排查 pprof")
	ix.Add(101, "channel 没人读 goroutine 泄漏")
	ix.Add(103, "显存估算 7b 模型 int4 量化")
	ix.Add(105, "goroutine 和线程的区别")
	ix.Add(102, "sqlite fts5 中文分词")

	hits := ix.Search("goroutine 泄漏", 10)
	if len(hits) < 2 {
		t.Fatalf("hits = %+v", hits)
	}
	if hits[0].TID != 101 {
		t.Fatalf("top = %+v", hits[0])
	}
	seen := map[int]bool{}
	for i, h := range hits {
		if seen[h.TID] {
			t.Fatalf("tid %d duplicated", h.TID)
		}
		seen[h.TID] = true
		if h.Score <= 0 {
			t.Fatalf("non-positive score %+v", h)
		}
		if i > 0 && h.Score > hits[i-1].Score {
			t.Fatal("not sorted")
		}
	}
	if !seen[105] || seen[102] {
		t.Fatalf("expected 105 hit and 102 not: %+v", hits)
	}
	if got := ix.Search("显存", 1); len(got) != 1 || got[0].TID != 103 {
		t.Fatalf("显存 → %+v", got)
	}
	if got := ix.Search("rust", 5); len(got) != 0 {
		t.Fatalf("rust → %+v", got)
	}
}
