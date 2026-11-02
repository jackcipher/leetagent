# 提示

1. 取 `strings.Index(out, "{")` 到 `strings.LastIndex(out, "}")`，这一步就顺带处理了代码块。
2. 先解到 `struct{ Score float64; Reason string }`，再检查 `Score == math.Trunc(Score)` 和范围；`"4"` 这种字符串会让 Unmarshal 报错，正好拒绝。
3. Judge 循环两次；模型调用本身的错误直接返回，不重试。
