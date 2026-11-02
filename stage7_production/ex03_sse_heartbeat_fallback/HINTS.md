# 提示

1. WriteSSE：用 strings.Builder 拼好再一次性写出；`strings.Split("", "\n")` 返回 `[""]`，正好写出 "data: \n"。
2. Heartbeat.Run：for + select { ctx.Done(); `_, ok := <-ticks` }；写入放进一个具名方法 `ping()`，里面 Lock/defer Unlock。
3. 降级：定义 `type countingSink struct{ inner TokenSink; n int }`，它的 Token 方法计数后转发。primary 失败且 n==0 才调 backup。
