//go:build solution

package ex05_hybrid_rrf

import "leetagent/stage3_rag/ex05_hybrid_rrf/solution"

type (
	Doc    = solution.Doc
	Ranked = solution.Ranked
	BM25   = solution.BM25
)

var (
	Tokenize = solution.Tokenize
	NewBM25  = solution.NewBM25
	RRF      = solution.RRF
)
