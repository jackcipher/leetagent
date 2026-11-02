// Package solution 是 m2_research_index 的参考答案。
package solution

import (
	"hash/fnv"
	"math"
	"sort"
	"strings"
	"unicode"
)

func isASCIIWord(r rune) bool {
	return r < 128 && (unicode.IsLetter(r) || unicode.IsDigit(r))
}

func flushHan(run []rune, out []string) []string {
	if len(run) == 1 {
		return append(out, string(run))
	}
	for i := 0; i+1 < len(run); i++ {
		out = append(out, string(run[i:i+2]))
	}
	return out
}

func Tokenize(s string) []string {
	var out []string
	var word []rune
	var han []rune
	for _, r := range strings.ToLower(s) {
		switch {
		case isASCIIWord(r):
			if len(han) > 0 {
				out = flushHan(han, out)
				han = han[:0]
			}
			word = append(word, r)
		case unicode.Is(unicode.Han, r):
			if len(word) > 0 {
				out = append(out, string(word))
				word = word[:0]
			}
			han = append(han, r)
		default:
			if len(word) > 0 {
				out = append(out, string(word))
				word = word[:0]
			}
			if len(han) > 0 {
				out = flushHan(han, out)
				han = han[:0]
			}
		}
	}
	if len(word) > 0 {
		out = append(out, string(word))
	}
	if len(han) > 0 {
		out = flushHan(han, out)
	}
	return out
}

type HashEmbedder struct {
	Dim int
}

func (e HashEmbedder) Embed(text string) []float64 {
	v := make([]float64, e.Dim)
	for _, tok := range Tokenize(text) {
		h := fnv.New32a()
		h.Write([]byte(tok))
		v[h.Sum32()%uint32(e.Dim)]++
	}
	var norm float64
	for _, x := range v {
		norm += x * x
	}
	if norm == 0 {
		return v
	}
	norm = math.Sqrt(norm)
	for i := range v {
		v[i] /= norm
	}
	return v
}

type Hit struct {
	TID   int
	Score float64
	Text  string
}

type entry struct {
	tid  int
	text string
	vec  []float64
}

type Index struct {
	Emb     HashEmbedder
	entries []entry
}

func NewIndex(dim int) *Index { return &Index{Emb: HashEmbedder{Dim: dim}} }

func (ix *Index) Add(tid int, text string) {
	ix.entries = append(ix.entries, entry{tid: tid, text: text, vec: ix.Emb.Embed(text)})
}

func dot(a, b []float64) float64 {
	var s float64
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}

type byScore []Hit

func (h byScore) Len() int      { return len(h) }
func (h byScore) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h byScore) Less(i, j int) bool {
	if h[i].Score != h[j].Score {
		return h[i].Score > h[j].Score
	}
	return h[i].TID < h[j].TID
}

func (ix *Index) Search(query string, k int) []Hit {
	q := ix.Emb.Embed(query)
	best := map[int]Hit{}
	for _, e := range ix.entries {
		s := dot(q, e.vec)
		if s <= 0 {
			continue
		}
		if old, ok := best[e.tid]; !ok || s > old.Score {
			best[e.tid] = Hit{TID: e.tid, Score: s, Text: e.text}
		}
	}
	hits := make([]Hit, 0, len(best))
	for _, h := range best {
		hits = append(hits, h)
	}
	sort.Sort(byScore(hits))
	if k >= 0 && len(hits) > k {
		hits = hits[:k]
	}
	return hits
}
