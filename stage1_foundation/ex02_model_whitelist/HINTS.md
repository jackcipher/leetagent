# 提示

1. `strings.Split(s, ",")` 对空串会返回 `[""]`，所以要在循环里 `TrimSpace` 后跳过空项。
2. 用一个 `map[string]bool` 同时完成"去重"和"判断默认值是否在列表里"。
3. `Resolve` 先 `TrimSpace`，再判断空串；返回错误用 `fmt.Errorf("%w: %s", ErrModelNotAllowed, m)`。
