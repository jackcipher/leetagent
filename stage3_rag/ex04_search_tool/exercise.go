//go:build !solution

package ex04_search_tool

import (
	"context"
	"errors"

	"leetagent/common"
)

var ErrEmptyQuery = errors.New("query is empty")

// Hit 是一条检索结果。
type Hit struct {
	ID    string  `json:"id"`
	Score float64 `json:"score"`
	Text  string  `json:"text"`
}

// Embedder 把文本变成向量。
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float64, error)
}

// Searcher 是向量索引（ex03 的 Index 就满足这个接口的形状）。
type Searcher interface {
	Search(q []float64, k int) ([]Hit, error)
}

// SearchPosts 是给模型用的 search_posts 工具。
type SearchPosts struct {
	Embedder Embedder
	Index    Searcher
	DefaultK int // top_k 缺省值
	MaxK     int // top_k 上限
}

// Def 返回工具定义：
//
//	name: "search_posts"
//	description: 非空，说明"在已入库的论坛帖子里做语义检索，返回最相关的片段和帖子 id"之类
//	parameters: JSON Schema object，properties: query(string, 必填), top_k(integer)；required: ["query"]
func (s *SearchPosts) Def() common.Tool {
	// TODO
	return common.Tool{}
}

// Run 解析 {"query": "...", "top_k": n}：
//
//	query 去空白后为空 -> ErrEmptyQuery；JSON 非法 -> 返回错误；
//	top_k <= 0 用 DefaultK，超过 MaxK 截到 MaxK；
//	Embed(query) -> Search(vec, k) -> 结果 JSON 数组（[]Hit 的 json.Marshal），没有结果返回 "[]"。
func (s *SearchPosts) Run(ctx context.Context, argsJSON string) (string, error) {
	// TODO
	return "", nil
}
