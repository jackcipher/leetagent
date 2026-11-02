# ex03 内存向量索引 + Top-K（Medium）

## 题目
在接入 sqlite-vec / pgvector 之前，先自己写一个最朴素的向量库：全部放内存，查询时暴力算一遍余弦，取前 K。
几万条以内，暴力扫描完全够用，而且是检验"向量库返回对不对"的基准。

## 接口
```go
func NewIndex() *Index
func (x *Index) Add(id string, vec []float64, text string) error
func (x *Index) Remove(id string) bool
func (x *Index) Len() int
func (x *Index) Search(q []float64, k int) ([]Hit, error)
```

## 要求
- 第一条向量决定维度；之后不一致返回 `ErrDimMismatch`。
- 同 ID 覆盖，不增加条数。
- 排序：分数降序，同分按 ID 升序（**结果必须确定**，否则测试和评测都会抖）。
- 建议存入时就单位化，查询时只做点积。

## 验收
`go test ./stage3_rag/ex03_vector_index/`

## 延伸
- 向量多了以后用 ANN（HNSW、IVF）近似检索；sqlite-vec 适合单机，pgvector 适合已有 Postgres。
- Top-K 也可以用大小为 K 的最小堆做到 O(n log k)。
