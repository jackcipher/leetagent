// Package solution 是 ex04_search_tool 的参考答案。
package solution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"leetagent/common"
)

var ErrEmptyQuery = errors.New("query is empty")

type Hit struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
	Text  string  `json:"text"`
}

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float64, error)
}

type Searcher interface {
	Search(q []float64, k int) ([]Hit, error)
}

type SearchPosts struct {
	Embedder Embedder
	Index    Searcher
	DefaultK int
	MaxK     int
}

const searchSchema = `{
  "type": "object",
  "properties": {
    "query": {"type": "string", "description": "要检索的问题或关键词，用自然语言描述"},
    "top_k": {"type": "integer", "description": "返回条数，默认 5"}
  },
  "required": ["query"]
}`

func (s *SearchPosts) Def() common.Tool {
	return common.NewFunctionTool("search_posts",
		"在已入库的论坛帖子里做语义检索，返回最相关的片段及其帖子 id。需要引用帖子内容回答问题时先调用它。",
		searchSchema)
}

type searchArgs struct {
	Query string `json:"query"`
	TopK  int    `json:"top_k"`
}

func (s *SearchPosts) Run(ctx context.Context, argsJSON string) (string, error) {
	var a searchArgs
	if err := json.Unmarshal([]byte(argsJSON), &a); err != nil {
		return "", fmt.Errorf("bad arguments: %w", err)
	}
	a.Query = strings.TrimSpace(a.Query)
	if a.Query == "" {
		return "", ErrEmptyQuery
	}
	k := a.TopK
	if k <= 0 {
		k = s.DefaultK
	}
	if s.MaxK > 0 && k > s.MaxK {
		k = s.MaxK
	}
	vec, err := s.Embedder.Embed(ctx, a.Query)
	if err != nil {
		return "", err
	}
	hits, err := s.Index.Search(vec, k)
	if err != nil {
		return "", err
	}
	if hits == nil {
		hits = []Hit{}
	}
	b, err := json.Marshal(hits)
	return string(b), err
}
