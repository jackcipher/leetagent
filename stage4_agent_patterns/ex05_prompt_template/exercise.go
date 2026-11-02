//go:build !solution

package ex05_prompt_template

import "leetagent/common"

// Section 是 system prompt 的一节。
type Section struct {
	Title string
	Body  string
}

// Build 把各节拼成 system prompt：每节渲染为 "## <Title>\n<Body 去首尾空白>"，节之间用一个空行分隔。
// Body 去空白后为空的节跳过。全部为空返回 ""。
func Build(sections []Section) string {
	// TODO
	return ""
}

// RenderToolGuide 生成"可用工具"一节的正文：按工具名字典序，每行 "- <name>: <description>"。
func RenderToolGuide(tools []common.Tool) string {
	// TODO
	return ""
}

// LintToolDescription 检查工具描述质量，返回问题列表（无问题返回空）。按以下顺序检查并追加：
//
//	"name: must be snake_case"      名字不匹配 ^[a-z][a-z0-9_]*$
//	"description: too short"        描述少于 10 个 rune
//	"description: too long"         描述多于 200 个 rune
//	"parameters: invalid json"      Parameters 不是合法 JSON 对象
//	"parameters.<p>: missing description"   properties 里每个缺少 description 的参数（按参数名排序）
func LintToolDescription(t common.Tool) []string {
	// TODO
	return nil
}
