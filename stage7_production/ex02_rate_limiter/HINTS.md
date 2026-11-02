# 提示

1. Bucket 字段：`capacity, rate, tokens float64; last time.Time; clock Clock`。NewBucket 时 tokens = capacity, last = clock.Now()。
2. Allow：`tokens += now.Sub(last).Seconds() * rate`，cap 到 capacity，更新 last，`tokens >= 1` 则减 1。
3. DailyQuota 存 `day string`（`Now().UTC().Format("2006-01-02")`）和 `used map[string]int`；每次操作前比较日期，不同就重建 map。中间件用 `http.Error(w, "rate limited", 429)`（它会自动加 "\n"），记得在之前设置 Retry-After。
