// Package solution 是 ex05_hybrid_rrf 的参考答案。
package solution

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

type Doc struct {
	ID   string
	Text string
}

type Ranked struct {
	ID    string
	Score float64
}

func isWordRune(r rune) bool {
	return r < 128 && (unicode.IsLetter(r) || unicode.IsDigit(r))
}

func Tokenize(s string) []string {
	var out []string
	var word strings.Builder
	for _, r := range s {
		if isWordRune(r) {
			word.WriteRune(unicode.ToLower(r))
			continue
		}
		if word.Len() > 0 {
			out = append(out, word.String())
			word.Reset()
		}
		if unicode.Is(unicode.Han, r) {
			out = append(out, string(r))
		}
	}
	if word.Len() > 0 {
		out = append(out, word.String())
	}
	return out
}

type docStat struct {
	id  string
	tf  map[string]int
	len int
}

type BM25 struct {
	K1, B  float64
	docs   []docStat
	df     map[string]int
	avgLen float64
}

func NewBM25(docs []Doc) *BM25 {
	b := &BM25{K1: 1.2, B: 0.75, df: map[string]int{}}
	total := 0
	for _, d := range docs {
		toks := Tokenize(d.Text)
		st := docStat{id: d.ID, tf: map[string]int{}, len: len(toks)}
		for _, t := range toks {
			st.tf[t]++
		}
		for t := range st.tf {
			b.df[t]++
		}
		total += len(toks)
		b.docs = append(b.docs, st)
	}
	if len(docs) > 0 {
		b.avgLen = float64(total) / float64(len(docs))
	}
	return b
}

type byScore []Ranked

func (r byScore) Len() int      { return len(r) }
func (r byScore) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r byScore) Less(i, j int) bool {
	if r[i].Score != r[j].Score {
		return r[i].Score > r[j].Score
	}
	return r[i].ID < r[j].ID
}

func uniq(toks []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, t := range toks {
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

func (b *BM25) Search(query string, k int) []Ranked {
	terms := uniq(Tokenize(query))
	n := float64(len(b.docs))
	var out []Ranked
	for _, d := range b.docs {
		score := 0.0
		for _, t := range terms {
			tf := float64(d.tf[t])
			if tf == 0 {
				continue
			}
			df := float64(b.df[t])
			idf := math.Log(1 + (n-df+0.5)/(df+0.5))
			norm := b.K1 * (1 - b.B + b.B*float64(d.len)/b.avgLen)
			score += idf * tf * (b.K1 + 1) / (tf + norm)
		}
		if score > 0 {
			out = append(out, Ranked{ID: d.id, Score: score})
		}
	}
	sort.Sort(byScore(out))
	if k >= 0 && len(out) > k {
		out = out[:k]
	}
	return out
}

func RRF(lists [][]string, k int) []Ranked {
	scores := map[string]float64{}
	for _, l := range lists {
		for i, id := range l {
			scores[id] += 1 / float64(k+i+1)
		}
	}
	out := make([]Ranked, 0, len(scores))
	for id, s := range scores {
		out = append(out, Ranked{ID: id, Score: s})
	}
	sort.Sort(byScore(out))
	return out
}
