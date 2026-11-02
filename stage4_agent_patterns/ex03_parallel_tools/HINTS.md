# 提示

1. 预分配 `results := make([]Result, len(calls))`，每个 goroutine 只写自己的 `&results[i]`——顺序天然保持，也没有竞争。
2. 信号量：`sem := make(chan struct{}, n)`，启动前 `sem <- struct{}{}`（和 `<-ctx.Done()` 一起 select），worker 结束时 `<-sem`。
3. 定义 `type job struct{ctx; runner; call; out *Result; sem chan struct{}; wg *sync.WaitGroup}` 和方法 `func (j job) run()`，里面先 `defer j.wg.Done()`。拿到信号量后再检查一次 `ctx.Err()`。
