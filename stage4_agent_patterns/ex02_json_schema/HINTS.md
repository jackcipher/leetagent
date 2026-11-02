# 提示

1. 把 schema 解析成递归 struct：`Properties map[string]*schemaNode`、`Items *schemaNode`，可选数字用指针（`*float64`）区分"没写"和"写了 0"；`AdditionalProperties *bool`。
2. 写一个递归函数 `check(s *schemaNode, v any, path string, errs []string) []string`：先查类型（不符直接 return），再查 enum，再 `switch x := v.(type)` 分 object/array/string/json.Number 处理。
3. integer 判断：`f, _ := n.Float64(); f == math.Trunc(f)`。enum 比较时把 json.Number 转成 float64，再和 schema 里的值（Unmarshal 出来的 float64/string/bool）比较。
