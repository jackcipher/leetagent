# 提示

1. LoadCases：bufio.Scanner 逐行，自己维护 line 计数（空行/注释行也要计数）；错误用 `fmt.Errorf("line %d: %w", line, err)`。
2. Grade：先把 out.Text 转小写，再和每个关键词的小写比较。三类检查顺序固定。
3. Markdown：百分比用 `%.1f%%`，PassRate()*100。注意 ✅ 那行原因列是两个空格 `|  |`。
