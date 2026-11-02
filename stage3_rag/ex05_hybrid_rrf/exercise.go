//go:build !solution

package ex05_hybrid_rrf

// Doc 是一篇待检索文档。
type Doc struct {
	ID   string
	Text string
}

// Ranked 是一条排序结果。
type Ranked struct {
	ID    string
	Score float64
}

// Tokenize 简易分词：
//
//	连续的 ASCII 字母/数字组成一个词并转小写；
//	每个汉字（unicode.Han）单独成一个词；
//	其它字符（空白、标点）作为分隔符丢弃。
func Tokenize(s string) []string {
	// TODO
	return nil
}

// BM25 是关键词检索器。
type BM25 struct {
	K1, B float64
	// TODO: 你需要的统计量
}

// NewBM25 用 K1=1.2, B=0.75 建立索引。
func NewBM25(docs []Doc) *BM25 {
	// TODO
	return &BM25{}
}

// Search 计算每篇文档对 query 的 BM25 分数，只返回分数 > 0 的，按分数降序、同分 ID 升序，最多 k 条。
//
//	idf(t)  = ln(1 + (N - df + 0.5) / (df + 0.5))
//	score   = Σ_t idf(t) * tf*(K1+1) / (tf + K1*(1 - B + B*len/avgLen))
//
// query 里重复的词只算一次。
func (b *BM25) Search(query string, k int) []Ranked {
	// TODO
	return nil
}

// RRF（Reciprocal Rank Fusion）融合多路排序：每路里排第 r 名（从 1 开始）的 id 得分 1/(k + r)，各路相加。
// 返回按分数降序、同分 ID 升序的结果。
func RRF(lists [][]string, k int) []Ranked {
	// TODO
	return nil
}
