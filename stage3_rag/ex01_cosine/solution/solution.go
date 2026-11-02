// Package solution 是 ex01_cosine 的参考答案。
package solution

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

var (
	ErrDimMismatch = errors.New("vector dimension mismatch")
	ErrZeroVector  = errors.New("zero vector")
	ErrBadResponse = errors.New("bad embedding response")
)

func Dot(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, ErrDimMismatch
	}
	s := 0.0
	for i := range a {
		s += a[i] * b[i]
	}
	return s, nil
}

func Norm(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x * x
	}
	return math.Sqrt(s)
}

func Cosine(a, b []float64) (float64, error) {
	d, err := Dot(a, b)
	if err != nil {
		return 0, err
	}
	na, nb := Norm(a), Norm(b)
	if na == 0 || nb == 0 {
		return 0, ErrZeroVector
	}
	return d / (na * nb), nil
}

func Normalize(v []float64) ([]float64, error) {
	n := Norm(v)
	if n == 0 {
		return nil, ErrZeroVector
	}
	out := make([]float64, len(v))
	for i, x := range v {
		out[i] = x / n
	}
	return out, nil
}

type embeddingResponse struct {
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func ParseEmbeddingResponse(body []byte) ([][]float64, error) {
	var r embeddingResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadResponse, err)
	}
	if len(r.Data) == 0 {
		return nil, fmt.Errorf("%w: empty data", ErrBadResponse)
	}
	out := make([][]float64, len(r.Data))
	for _, d := range r.Data {
		if d.Index < 0 || d.Index >= len(out) || out[d.Index] != nil {
			return nil, fmt.Errorf("%w: bad index %d", ErrBadResponse, d.Index)
		}
		out[d.Index] = d.Embedding
	}
	return out, nil
}
