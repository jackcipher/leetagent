# 提示

1. `rs := []rune(s)`，之后 `len(rs)` 就是字符数，`string(rs[:n])` 一定是合法 UTF-8。
2. 省略数 = `len(rs) - maxRunes`。
3. `TruncateToolResult` 用 `max, ok := limits[tool]`，`!ok` 时用 defaultMax，然后直接调 `Truncate`。
