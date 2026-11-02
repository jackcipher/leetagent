// Package solution 是 ex02_chunking 的参考答案。
package solution

import (
	"errors"
	"strings"
)

var ErrBadParams = errors.New("size must be > 0 and 0 <= overlap < size")

type Chunk struct {
	Index int
	Start int
	Text  string
}

func SplitFixed(text string, size, overlap int) ([]Chunk, error) {
	if size <= 0 || overlap < 0 || overlap >= size {
		return nil, ErrBadParams
	}
	rs := []rune(text)
	var out []Chunk
	step := size - overlap
	for start := 0; start < len(rs); start += step {
		end := start + size
		if end > len(rs) {
			end = len(rs)
		}
		out = append(out, Chunk{Index: len(out), Start: start, Text: string(rs[start:end])})
		if end == len(rs) {
			break
		}
	}
	return out, nil
}

func SplitParagraphs(text string, size, overlap int) ([]Chunk, error) {
	if size <= 0 || overlap < 0 || overlap >= size {
		return nil, ErrBadParams
	}
	var paras []string
	for _, p := range strings.Split(text, "\n\n") {
		p = strings.TrimSpace(p)
		if p != "" {
			paras = append(paras, p)
		}
	}
	var out []Chunk
	cur := ""
	for _, p := range paras {
		if runeLen(p) > size {
			out = flush(out, cur)
			cur = ""
			parts, _ := SplitFixed(p, size, overlap)
			for _, c := range parts {
				out = flush(out, c.Text)
			}
			continue
		}
		if cur == "" {
			cur = p
		} else if runeLen(cur)+1+runeLen(p) <= size {
			cur = cur + "\n" + p
		} else {
			out = flush(out, cur)
			cur = p
		}
	}
	out = flush(out, cur)
	return out, nil
}

func runeLen(s string) int { return len([]rune(s)) }

func flush(out []Chunk, text string) []Chunk {
	if text == "" {
		return out
	}
	return append(out, Chunk{Index: len(out), Start: -1, Text: text})
}
