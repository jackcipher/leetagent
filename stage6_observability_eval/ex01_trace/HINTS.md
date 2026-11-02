# 提示

1. Span 需要知道时钟才能 End()：加一个未导出字段 `tracer *Tracer`，在 Start 时赋值。
2. Render 写一个递归函数 `render(b *strings.Builder, s *Span, depth int)`：缩进 `strings.Repeat("  ", depth)`，耗时 `s.Duration().Milliseconds()`。
3. 属性先收集 key 再 sort.Strings；Find 也是递归先序遍历。
