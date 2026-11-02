// Package solution 是 ex03_vector_index 的参考答案。
package solution

import (
	"errors"
	"math"
	"sort"
)

var (
	ErrDimMismatch = errors.New("vector dimension mismatch")
	ErrEmptyID     = errors.New("empty id")
)

type Hit struct {
	ID    string
	Score float64
	Text  string
}

type item struct {
	vec  []float64 // 已单位化（零向量保持全 0）
	text string
}

type Index struct {
	dim   int
	items map[string]item
}

func NewIndex() *Index { return &Index{items: map[string]item{}} }

func unit(v []float64) []float64 {
	s := 0.0
	for _, x := range v {
		s += x * x
	}
	out := make([]float64, len(v))
	if s == 0 {
		return out
	}
	n := math.Sqrt(s)
	for i, x := range v {
		out[i] = x / n
	}
	return out
}

func (x *Index) Add(id string, vec []float64, text string) error {
	if id == "" {
		return ErrEmptyID
	}
	if len(x.items) == 0 && x.dim == 0 {
		x.dim = len(vec)
	}
	if len(vec) != x.dim {
		return ErrDimMismatch
	}
	x.items[id] = item{vec: unit(vec), text: text}
	return nil
}

func (x *Index) Remove(id string) bool {
	_, ok := x.items[id]
	delete(x.items, id)
	return ok
}

func (x *Index) Len() int { return len(x.items) }

// byScore 实现 sort.Interface：分数降序，同分按 ID 升序。
type byScore []Hit

func (h byScore) Len() int      { return len(h) }
func (h byScore) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h byScore) Less(i, j int) bool {
	if h[i].Score != h[j].Score {
		return h[i].Score > h[j].Score
	}
	return h[i].ID < h[j].ID
}

func (x *Index) Search(q []float64, k int) ([]Hit, error) {
	if k <= 0 || len(x.items) == 0 {
		return nil, nil
	}
	if len(q) != x.dim {
		return nil, ErrDimMismatch
	}
	uq := unit(q)
	hits := make([]Hit, 0, len(x.items))
	for id, it := range x.items {
		s := 0.0
		for i := range uq {
			s += uq[i] * it.vec[i]
		}
		hits = append(hits, Hit{ID: id, Score: s, Text: it.text})
	}
	sort.Sort(byScore(hits))
	if len(hits) > k {
		hits = hits[:k]
	}
	return hits, nil
}
