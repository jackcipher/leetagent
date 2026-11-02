# ex05 Prompt 模板与工具描述体检（Easy）

## 题目
Prompt 是 Agent 的"源代码"。把 system prompt 结构化成若干 Section，可以版本化、单测、做 A/B。
工具描述写得差，小模型就会乱调工具——所以给工具描述写一个 linter。

## 接口
```go
type Section struct{ Title, Body string }
func Build(sections []Section) string
func RenderToolGuide(tools []common.Tool) string
func LintToolDescription(t common.Tool) []string
```
细节见 exercise.go 注释。

## 写好工具描述的经验
- 名字用 snake_case 动词开头：`get_post`、`search_posts`
- 描述写清楚"什么时候用"和"返回什么"，10~200 字
- 每个参数都要有 description，最好带示例和单位
- 用 本地 LLM 测试：同一个问题，改描述前后模型选择工具的正确率差多少？

## 验收
`go test ./stage4_agent_patterns/ex05_prompt_template/`
