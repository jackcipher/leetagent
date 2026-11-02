//go:build !solution

package ex02_chunking

import "errors"

var ErrBadParams = errors.New("size must be > 0 and 0 <= overlap < size")

// Chunk 是切出来的一块。Start 是它在原文中的起始 rune 下标。
type Chunk struct {
	Index int
	Start int
	Text  string
}

// SplitFixed 按 rune 定长切块，相邻块重叠 overlap 个字符：
//
//	第 i 块从 i*(size-overlap) 开始，长度 size（最后一块可能更短）；
//	当某块已经到达文本末尾就停止（不要产生完全被上一块包含的尾块）。
//
// 空文本返回空切片（nil 也可以）。参数非法返回 ErrBadParams。
func SplitFixed(text string, size, overlap int) ([]Chunk, error) {
	// TODO
	return nil, nil
}

// SplitParagraphs 先按空行（"\n\n"）分段，去掉每段首尾空白并丢弃空段，
// 然后把相邻段落贪心合并：合并后（段间用 "\n" 连接）不超过 size 就合并；
// 单段超过 size 时，用 SplitFixed(段, size, overlap) 切开，各块单独成 chunk。
// 结果 Index 从 0 递增；Start 统一填 -1（段落模式不追踪位置）。
func SplitParagraphs(text string, size, overlap int) ([]Chunk, error) {
	// TODO
	return nil, nil
}
