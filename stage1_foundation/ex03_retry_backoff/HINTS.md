# 提示

1. `errors.Is(err, context.Canceled)` 能穿透 `%w` 包装；`errors.As(err, &se)`（`var se *HTTPStatusError`）能取出状态码错误。
2. `Delay` 用循环每次乘 2，一旦 >= MaxDelay 就返回 MaxDelay，避免溢出。
3. `Do` 的循环：`for i := 1; i <= attempts; i++ { err = op.Run(ctx); if err == nil || !IsRetryable(err) { return err }; if i == attempts { break }; if e := s.Sleep(ctx, p.Delay(i)); e != nil { return e } }`。
