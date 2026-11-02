# 提示

1. 字段只需要一个 `map[string]Tool`。别忘了在 `NewRegistry` 里 `make`，否则写 nil map 会 panic。
2. map 遍历顺序是随机的。先把 key 收集到 `[]string`，`sort.Strings`，再按顺序取 `Def()`。
3. 错误包装：`fmt.Errorf("%w: %s", ErrUnknownTool, name)`，这样 `errors.Is` 能识别，信息里也有名字。
