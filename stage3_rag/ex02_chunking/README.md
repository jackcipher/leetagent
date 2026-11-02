# ex02 切块（Chunking）与重叠（Medium）

## 题目
一个帖子太长，整篇做 embedding 会把多个主题"平均"掉，检索不准；也塞不进上下文。
所以入库前要切块。实现两种切法：
1. **定长 + 重叠**：简单稳定，重叠保证跨块的句子不会被腰斩后完全丢失。
2. **按段落合并**：尊重文本结构，块内语义更完整。

## 接口
```go
func SplitFixed(text string, size, overlap int) ([]Chunk, error)
func SplitParagraphs(text string, size, overlap int) ([]Chunk, error)
```

## 图示（size=4, overlap=2）
```
一 二 三 四 五 六 七
[一 二 三 四]
      [三 四 五 六]
            [五 六 七]
```

## 验收
`go test ./stage3_rag/ex02_chunking/`

## 经验值
- 中文帖子常用 300~500 字一块、50~100 字重叠。
- 每块保留元数据（帖子 tid、标题、楼层），检索命中后才能"带出处"回答（毕业项目要用）。
