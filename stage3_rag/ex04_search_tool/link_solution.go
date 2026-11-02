//go:build solution

package ex04_search_tool

import "leetagent/stage3_rag/ex04_search_tool/solution"

type (
	Hit         = solution.Hit
	Embedder    = solution.Embedder
	Searcher    = solution.Searcher
	SearchPosts = solution.SearchPosts
)

var ErrEmptyQuery = solution.ErrEmptyQuery
