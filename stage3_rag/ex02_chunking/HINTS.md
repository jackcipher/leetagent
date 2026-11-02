# 提示

1. 步长 `step := size - overlap`；`for start := 0; start < len(rs); start += step`，每次 `end = min(start+size, len(rs))`。
2. `end == len(rs)` 时 break，否则 "abcdef" 会多出一个 "ef" 尾块。
3. 段落合并：维护 `cur`，能拼就拼（注意算上 "\n" 的 1 个字符），拼不下就把 cur 输出再换新；超长段先输出 cur，再把 SplitFixed 的每块逐个输出。
