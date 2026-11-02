# 提示

1. `slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))`。
2. 动态选择级别用 `l.LogAttrs(ctx, level, "agent round", attrs...)`，attrs 是 `[]slog.Attr`。
3. 毫秒：`d.Milliseconds()` 返回 int64，用 `slog.Int64`。
