//go:build !solution

// Package m2_research_index 是毕业项目第 2 步：离线可用的向量索引（不需要 embedding 模型）。
package m2_research_index

// Tokenize 把文本切成词：
//   - 先转小写；
//   - 连续的 ASCII 字母/数字组成一个词（如 "goroutine"、"7b"）；
//   - 连续的汉字（unicode.Han）按**相邻两字**切成 bigram，"显存估算" → "显存","存估","估算"；只有一个汉字时就是它本身；
//   - 其它字符（空格、标点）都是分隔符。
func Tokenize(s string) []string {
	// TODO
	return nil
}

// HashEmbedder 用"特征哈希"把文本变成 Dim 维向量：每个 token 用 fnv-32a 哈希，下标 = hash % Dim，该维 +1；
// 最后做 L2 归一化（全零向量保持全零）。
type HashEmbedder struct {
	Dim int
}

func (e HashEmbedder) Embed(text string) []float64 {
	// TODO
	return nil
}

// Hit 是一条检索结果。
type Hit struct {
	TID   int
	Score float64
	Text  string
}

// Index 存放多个 chunk，同一个帖子（tid）可以有多个 chunk。
type Index struct {
	Emb HashEmbedder
	// TODO: 你需要的字段
}

// NewIndex 创建索引。
func NewIndex(dim int) *Index {
	// TODO
	return nil
}

// Add 加入一个 chunk。
func (ix *Index) Add(tid int, text string) {
	// TODO
}

// Search 用余弦相似度（向量已归一化 = 点积）检索：
// 同一个 tid 只保留得分最高的 chunk；去掉 Score<=0 的；按 Score 降序（相同按 TID 升序）取前 k 条。
func (ix *Index) Search(query string, k int) []Hit {
	// TODO
	return nil
}
