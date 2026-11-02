//go:build solution

package ex02_chunking

import "leetagent/stage3_rag/ex02_chunking/solution"

type Chunk = solution.Chunk

var (
	ErrBadParams    = solution.ErrBadParams
	SplitFixed      = solution.SplitFixed
	SplitParagraphs = solution.SplitParagraphs
)
