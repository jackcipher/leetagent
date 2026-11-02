# 提示

1. Detect：文本转小写，每条规则任一短语命中就记录规则名并 break（每条最多一次）。
2. 大小写不敏感替换用 `regexp.MustCompile("(?i)</untrusted")` 的 ReplaceAllString。
3. GuardedRunner.Call 顺序：白名单 → 执行 → Wrap → Detect（对原始 out 检测）→ 加警告并记录 Alerts。
