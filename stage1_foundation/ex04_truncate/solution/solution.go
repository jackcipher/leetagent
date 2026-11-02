// Package solution 是 ex04_truncate 的参考答案。
package solution

import "fmt"

const TruncMarker = "\n…[已截断，省略 %d 字]"

const MiddleMarker = "\n…[中间省略 %d 字]…\n"

func Truncate(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return s
	}
	rs := []rune(s)
	if len(rs) <= maxRunes {
		return s
	}
	return string(rs[:maxRunes]) + fmt.Sprintf(TruncMarker, len(rs)-maxRunes)
}

func TruncateHeadTail(s string, head, tail int) string {
	rs := []rune(s)
	if head < 0 {
		head = 0
	}
	if tail < 0 {
		tail = 0
	}
	if len(rs) <= head+tail {
		return s
	}
	omitted := len(rs) - head - tail
	return string(rs[:head]) + fmt.Sprintf(MiddleMarker, omitted) + string(rs[len(rs)-tail:])
}

func TruncateToolResult(tool, result string, limits map[string]int, defaultMax int) string {
	max, ok := limits[tool]
	if !ok {
		max = defaultMax
	}
	return Truncate(result, max)
}
