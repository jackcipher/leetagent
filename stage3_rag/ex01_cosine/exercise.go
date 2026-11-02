//go:build !solution

package ex01_cosine

import "errors"

var (
	ErrDimMismatch = errors.New("vector dimension mismatch")
	ErrZeroVector  = errors.New("zero vector")
	ErrBadResponse = errors.New("bad embedding response")
)

// Dot 点积。长度不同返回 ErrDimMismatch。
func Dot(a, b []float64) (float64, error) {
	// TODO
	return 0, nil
}

// Norm L2 范数。
func Norm(v []float64) float64 {
	// TODO
	return 0
}

// Cosine 余弦相似度 = dot / (|a|*|b|)。长度不同 ErrDimMismatch，任一为零向量 ErrZeroVector。
func Cosine(a, b []float64) (float64, error) {
	// TODO
	return 0, nil
}

// Normalize 返回单位化后的新向量（不修改入参）。零向量返回 ErrZeroVector。
func Normalize(v []float64) ([]float64, error) {
	// TODO
	return nil, nil
}

// ParseEmbeddingResponse 解析 OpenAI /v1/embeddings 的响应体：
//
//	{"data":[{"index":1,"embedding":[...]},{"index":0,"embedding":[...]}], "model":"..."}
//
// 按 index 排序后返回向量列表。data 为空、index 不连续(不是 0..n-1)或 JSON 非法时返回 ErrBadResponse。
func ParseEmbeddingResponse(body []byte) ([][]float64, error) {
	// TODO
	return nil, nil
}
