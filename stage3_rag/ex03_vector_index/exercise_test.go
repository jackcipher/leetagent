package ex03_vector_index

import (
	"errors"
	"math"
	"testing"
)

func TestIndex(t *testing.T) {
	x := NewIndex()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(x.Add("go", []float64{1, 0, 0}, "Go 并发"))
	must(x.Add("rust", []float64{0, 1, 0}, "Rust 所有权"))
	must(x.Add("go2", []float64{2, 0.2, 0}, "Go 调度器"))
	must(x.Add("cook", []float64{0, 0, 3}, "红烧肉"))
	if x.Len() != 4 {
		t.Fatalf("Len = %d", x.Len())
	}
	if err := x.Add("bad", []float64{1, 2}, ""); !errors.Is(err, ErrDimMismatch) {
		t.Errorf("dim err = %v", err)
	}
	if err := x.Add("", []float64{1, 2, 3}, ""); !errors.Is(err, ErrEmptyID) {
		t.Errorf("empty id err = %v", err)
	}

	hits, err := x.Search([]float64{1, 0.05, 0}, 2)
	if err != nil || len(hits) != 2 {
		t.Fatalf("hits=%v err=%v", hits, err)
	}
	if hits[0].ID != "go" && hits[0].ID != "go2" {
		t.Fatalf("top = %s", hits[0].ID)
	}
	if hits[0].Score < hits[1].Score || hits[1].Text == "" {
		t.Fatalf("not sorted / no text: %+v", hits)
	}

	// 同分按 ID 升序
	must(x.Add("a_dup", []float64{0, 0, 1}, "dup"))
	hits, _ = x.Search([]float64{0, 0, 1}, 2)
	if hits[0].ID != "a_dup" || hits[1].ID != "cook" || math.Abs(hits[0].Score-1) > 1e-9 {
		t.Fatalf("tie order: %+v", hits)
	}

	// 覆盖
	must(x.Add("cook", []float64{1, 0, 0}, "改了"))
	if x.Len() != 5 {
		t.Fatalf("overwrite changed Len = %d", x.Len())
	}
	if !x.Remove("cook") || x.Remove("cook") {
		t.Fatal("Remove semantics wrong")
	}
	if _, err := x.Search([]float64{1}, 1); !errors.Is(err, ErrDimMismatch) {
		t.Errorf("query dim err = %v", err)
	}
	if h, _ := x.Search([]float64{1, 0, 0}, 0); len(h) != 0 {
		t.Error("k=0 should return empty")
	}
	if h, _ := x.Search([]float64{1, 0, 0}, 100); len(h) != 4 {
		t.Errorf("k>len should return all, got %d", len(h))
	}
}
