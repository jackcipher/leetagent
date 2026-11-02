# 提示

1. 不能改入参：先 `append(durations(nil), values...)` 复制；按规范不用闭包，定义 `type durations []time.Duration` 实现 sort.Interface。
2. `rank := int(math.Ceil(p / 100 * float64(n)))`，再夹在 [1, n] 之间，返回 `sorted[rank-1]`。
3. Summarize 一次遍历收集 ttft/total 切片、tokens 总和、生成耗时总和（End-FirstToken）和费用。
