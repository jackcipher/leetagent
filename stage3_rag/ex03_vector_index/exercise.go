//go:build !solution

package ex03_vector_index

import "errors"

var (
	ErrDimMismatch = errors.New("vector dimension mismatch")
	ErrEmptyID     = errors.New("empty id")
)

// Hit 是一条检索结果。
type Hit struct {
	ID    string
	Score float64 // 余弦相似度
	Text  string
}

// Index 是内存向量索引（暴力扫描）。第一条向量决定维度。
type Index struct {
	// TODO
}

// NewIndex 创建空索引。
func NewIndex() *Index {
	// TODO
	return &Index{}
}

// Add 加入或覆盖（同 ID）一条向量。ID 为空 ErrEmptyID；维度和已有不一致 ErrDimMismatch；零向量也允许存但永远得分 0。
func (x *Index) Add(id string, vec []float64, text string) error {
	// TODO
	return nil
}

// Remove 删除，返回是否存在。
func (x *Index) Remove(id string) bool {
	// TODO
	return false
}

// Len 返回条数。
func (x *Index) Len() int {
	// TODO
	return 0
}

// Search 返回与 q 余弦相似度最高的 k 条，按 Score 降序，分数相同按 ID 升序。
// k<=0 或索引为空返回空；q 维度不对返回 ErrDimMismatch。
func (x *Index) Search(q []float64, k int) ([]Hit, error) {
	// TODO
	return nil, nil
}
