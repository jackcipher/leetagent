//go:build !solution

package ex04_router

import (
	"context"

	"leetagent/common"
)

// Route 是一个可路由的目标（某个子 Agent 或处理流程）。
type Route struct {
	Name        string   // 小写英文，例如 "forum"、"code"、"chitchat"
	Description string   // 给 LLM 分类器看的说明
	Keywords    []string // 关键词（大小写不敏感的子串匹配）
}

// KeywordRoute 统计每个路由命中的关键词个数（同一个关键词只算一次），返回命中最多的路由名；
// 并列取 routes 中靠前的；全部为 0 返回 ("", false)。
func KeywordRoute(routes []Route, text string) (string, bool) {
	// TODO
	return "", false
}

// ClassifyPrompt 生成分类用的 system prompt，格式严格为：
//
//	"你是一个意图分类器。根据用户输入选择最合适的类别，只输出类别名，不要输出其它内容。\n类别：\n"
//	+ 每个路由一行 "- <name>: <description>\n"
func ClassifyPrompt(routes []Route) string {
	// TODO
	return ""
}

// Router 先用关键词路由；没命中再问 LLM。
type Router struct {
	LLM     common.Chatter
	Model   string
	Routes  []Route
	Default string
}

// Route 返回 (路由名, 来源)，来源为 "keyword" / "llm" / "default"。
//
//	LLM 调用：消息 [{system: ClassifyPrompt}, {user: text}]，不带工具。
//	解析 LLM 输出：去首尾空白、反引号、句号（。.）、转小写；若恰好等于某个路由名则采用；
//	  否则若输出中**恰好只包含一个**路由名作为子串，也采用；否则使用 Default。
//	LLM 出错 -> Default。
func (r *Router) Route(ctx context.Context, text string) (string, string) {
	// TODO
	return "", ""
}
