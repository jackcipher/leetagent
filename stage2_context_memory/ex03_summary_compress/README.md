# ex03 对话摘要压缩（Medium）

## 题目
ex02 的裁剪是"直接丢"，会丢信息。更好的做法：把旧对话交给模型总结成一段摘要，用摘要替换原文。
这就是 ChatGPT / Claude 长对话里"compaction"的基本形态。

## 接口
```go
func RenderTranscript(msgs []common.Message) string
func CountTurns(msgs []common.Message) int
func Compress(ctx context.Context, c common.Chatter, model string, msgs []common.Message, keepTurns int) ([]common.Message, error)
```

## 输出结构
```
[原 system] [system: "以下是之前对话的摘要：\n…"] [最近 keepTurns 轮原文]
```

## 要求
见 exercise.go 注释，关键点：
- 只有轮数 > keepTurns 才调模型（省钱）。
- 发给摘要模型的不要带 tools（避免它去调工具）。
- 渲染 transcript 时工具调用要可读：`assistant -> get_list({"page":1})`。

## 验收
`go test ./stage2_context_memory/ex03_summary_compress/`

## 思考
- 摘要也会"失真"。重要事实（用户名、偏好）应该走 ex04 的长期记忆，而不是指望摘要。
- 用小模型（Qwen 0.5B）做摘要、大模型做对话，是常见的降本搭配。
