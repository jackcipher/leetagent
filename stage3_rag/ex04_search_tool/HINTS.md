# 提示

1. 定义参数 struct：字段 `Query string` 带 json tag `query`，`TopK int` 带 json tag `top_k`；json.Unmarshal 后对 Query 做 TrimSpace。
2. Schema 用反引号多行字符串常量，`common.NewFunctionTool(name, desc, schema)` 直接构造。
3. `json.Marshal([]Hit(nil))` 输出 `null`；`json.Marshal([]Hit{})` 输出 `[]`。
