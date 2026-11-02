# ex04 工具结果截断（Easy）

## 题目
`fetch_one` 一个长帖可能几万字，原样回填会把 0.5B/Flash 模型的上下文撑爆，甚至直接报错。
实现三个截断函数，**按 rune 计数**，中文、emoji 都不能被切成乱码。

## 接口
```go
func Truncate(s string, maxRunes int) string
func TruncateHeadTail(s string, head, tail int) string
func TruncateToolResult(tool, result string, limits map[string]int, defaultMax int) string
```

## 要求
- `Truncate`：长度 <= max 原样；否则前 max 个 rune + `fmt.Sprintf(TruncMarker, 省略数)`；max<=0 不限制。
- `TruncateHeadTail`：保留头 head、尾 tail，中间换成 `MiddleMarker`。
- `TruncateToolResult`：每个工具可以有自己的上限。
- 告诉模型"被截断了"很重要：否则它会以为帖子就这么短，自信地编结论。

## 坑
`s[:n]` 是按**字节**切，一个汉字 3 字节，会切出 `\xe4\xbd` 这样的半个字。要用 `[]rune(s)`。

## 验收
`go test ./stage1_foundation/ex04_truncate/`
