# ex02 结构化输出与 JSON Schema 校验（Hard）

## 题目
Agent 经常需要模型输出**机器可读**的结果：评分、分类、提取字段、计划。
即使上游支持 `response_format: {type: "json_schema"}`，本地小模型也常常不遵守。
所以后端必须：**校验 → 把错误反馈给模型 → 重试**。

## 接口
```go
func Validate(schema, data json.RawMessage) []string
func ExtractJSON(output string) string
func AskStructured(ctx, c common.Chatter, model string, msgs []common.Message, schema json.RawMessage, retries int) (json.RawMessage, error)
```

## 支持的 Schema 子集
`type properties required additionalProperties(false) items enum minimum maximum minLength maxLength`

错误格式：`<path>: <原因>`，原因固定为以下之一（测试逐字比对）：
| 情况 | 原因文本 |
|---|---|
| 类型不对 | `expected <type>` |
| 缺字段 | `required`（path 是缺失字段的路径，如 `$.score`） |
| 多余字段 | `additional property not allowed` |
| 不在枚举 | `not in enum` |
| 数字越界 | `less than minimum <v>` / `greater than maximum <v>`（%v 格式） |
| 字符串长度 | `shorter than <n>` / `longer than <n>` |

## 坑
- `json.Unmarshal` 到 `any` 时数字都是 float64，没法区分 3 和 3.5 之外的大整数精度问题——用 `Decoder.UseNumber()`。
- 错误要排序，否则 map 遍历顺序会让测试随机失败。
- 重试时把**具体错误**告诉模型，比单纯说"格式错了"成功率高得多。

## 验收
`go test ./stage4_agent_patterns/ex02_json_schema/`
