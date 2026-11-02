# ex01 Embedding 与余弦相似度（Easy）

## 题目
RAG 的第一步：把文本变成向量（embedding），用余弦相似度衡量"意思有多近"。
实现向量基础运算，并解析 OpenAI 兼容 `/v1/embeddings` 的响应。

## 接口
```go
func Dot(a, b []float64) (float64, error)
func Norm(v []float64) float64
func Cosine(a, b []float64) (float64, error)
func Normalize(v []float64) ([]float64, error)
func ParseEmbeddingResponse(body []byte) ([][]float64, error)
```

## 协议
请求：`POST /v1/embeddings {"model":"...","input":["文本1","文本2"]}`
响应：`{"data":[{"index":0,"embedding":[...]}, ...]}` —— **data 不保证按 index 排序**，必须按 index 放回。

## 验收
`go test ./stage3_rag/ex01_cosine/`

## 知识点
- 向量先 Normalize，之后 cosine 就等于 dot，检索时省掉两次开方。
- 同一个库里只能用**同一个 embedding 模型**的向量，换模型要全量重建索引。
