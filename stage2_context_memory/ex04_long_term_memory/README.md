# ex04 长期记忆（Medium）

## 题目
摘要是"这次对话的记忆"，长期记忆是"跨对话的记忆"：用户叫什么、偏好中文、在学 Agent……
做法：每次用户发言后，让一个（可以是小）模型抽取事实 → 按 key 覆盖写入 → 下次对话把事实渲染进 system prompt。

## 接口
```go
func NewMemoryStore() *MemoryStore
func (m *MemoryStore) Upsert(f Fact)
func (m *MemoryStore) All() []Fact
func ParseFacts(output string) ([]Fact, error)
func (m *MemoryStore) Extract(ctx, c common.Chatter, model, userMsg string) ([]Fact, error)
func (m *MemoryStore) Render(max int) string
```

## 要求
- 模型经常把 JSON 包在 ```json 代码块里，`ParseFacts` 要能剥掉。
- 非 JSON 输出返回 `ErrBadJSON`（可 errors.Is）。`[]` 返回空切片（不是 nil）。
- `Render` 格式见注释，没事实返回空串（避免注入一个空标题浪费 token）。

## 验收
`go test ./stage2_context_memory/ex04_long_term_memory/`

## 延伸
- 生产里记忆存在 SQLite（参考 ex01），按 user_id 分区。
- 记忆多了以后不能全注入：用阶段 3 的向量检索"按需召回"。
- 安全：记忆来自用户输入，同样可能含注入指令（阶段 7）。
