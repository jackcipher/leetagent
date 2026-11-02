//go:build !solution

package ex04_truncate

// Truncate 保证返回值的"正文部分"不超过 maxRunes 个字符（按 rune，不按 byte，中文不能切坏）。
// 没超过时原样返回；超过时返回 前 maxRunes 个字符 + fmt.Sprintf(TruncMarker, 被省略的字符数)。
// maxRunes <= 0 表示不限制。
func Truncate(s string, maxRunes int) string {
	// TODO
	return ""
}

// TruncMarker 是截断后追加的提示，%d 是省略的字符数。模型看到它就知道内容不完整。
const TruncMarker = "\n…[已截断，省略 %d 字]"

// TruncateHeadTail 保留开头 head 个字符和结尾 tail 个字符，中间替换为 fmt.Sprintf(MiddleMarker, 省略数)。
// 若总长度 <= head+tail 原样返回。适合日志、报错堆栈这类"头尾都重要"的文本。
func TruncateHeadTail(s string, head, tail int) string {
	// TODO
	return ""
}

// MiddleMarker 是中间省略的提示。
const MiddleMarker = "\n…[中间省略 %d 字]…\n"

// TruncateToolResult 是给 Agent 循环用的入口：按工具名查 limits 里的上限，查不到用 defaultMax。
func TruncateToolResult(tool, result string, limits map[string]int, defaultMax int) string {
	// TODO
	return ""
}
