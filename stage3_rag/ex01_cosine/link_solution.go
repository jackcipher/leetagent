//go:build solution

package ex01_cosine

import "leetagent/stage3_rag/ex01_cosine/solution"

var (
	ErrDimMismatch         = solution.ErrDimMismatch
	ErrZeroVector          = solution.ErrZeroVector
	ErrBadResponse         = solution.ErrBadResponse
	Dot                    = solution.Dot
	Norm                   = solution.Norm
	Cosine                 = solution.Cosine
	Normalize              = solution.Normalize
	ParseEmbeddingResponse = solution.ParseEmbeddingResponse
)
