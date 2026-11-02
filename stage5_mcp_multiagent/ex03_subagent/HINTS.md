# 提示

1. Run 就是 stage4 ex01 的循环，只是消息由 System + task 新建。Tools 为 nil 时把调用当作 "unknown tool" 错误回填，不要 panic。
2. Call：先判断 `strings.HasPrefix(name, "delegate_")` 并在 Subs 里查 `strings.TrimPrefix(name, "delegate_")`；都满足才委派，否则走 Fallback。
3. Defs：收集 map 的 key，sort.Strings，再依次 DelegateTool。schema 可以写成一个字符串常量。
