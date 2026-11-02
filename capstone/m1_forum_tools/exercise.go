//go:build !solution

// Package m1_forum_tools 是毕业项目第 1 步：把论坛数据包装成 Agent 可调用的工具。
package m1_forum_tools

import (
	"context"

	"leetagent/common"
)

// Post 是一个帖子。
type Post struct {
	TID     int    `json:"tid"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Replies int    `json:"replies"`
	Content string `json:"content"`
}

// Summary 是列表/搜索里的简略信息（不含正文，省 token）。
type Summary struct {
	TID     int    `json:"tid"`
	Title   string `json:"title"`
	Replies int    `json:"replies"`
}

// Forum 是内存里的论坛。
type Forum struct {
	Posts []Post
}

// LoadForum 读取 JSON 数组文件。
func LoadForum(path string) (*Forum, error) {
	// TODO
	return nil, nil
}

// Defs 返回三个工具定义（顺序：get_list, get_post, search_posts），参数 schema 必须是合法 JSON：
//
//	get_list     {"limit": integer}            —— 按 replies 降序（相同按 tid 升序）返回前 limit 条 Summary；limit<=0 或缺省视为 10
//	get_post     {"tid": integer}  required tid —— 返回完整 Post
//	search_posts {"query": string} required query —— 标题或正文包含 query（大小写不敏感）的 Summary，按 tid 升序
func (f *Forum) Defs() []common.Tool {
	// TODO
	return nil
}

// Call 实现 common.ToolRunner，返回 JSON 字符串。错误（必须返回 error）：
//
//	未知工具   → "unknown tool: <name>"
//	参数不是 JSON → "bad arguments: ..."
//	帖子不存在 → "post <tid> not found"
//	query 为空 → "query is required"
//
// 结果为空列表时返回 "[]" 而不是 "null"。
func (f *Forum) Call(ctx context.Context, name, argsJSON string) (string, error) {
	// TODO
	return "", nil
}
