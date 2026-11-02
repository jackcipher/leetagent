# 提示

1. Gate.Call 就是三个分支：非危险直接调 Inner；危险先问 Approver；按结果记录 Log。`common.Call("", name, args)` 能直接造一个 ToolCall。
2. ChanApprover 字段：`mu sync.Mutex; seq int; reqs chan Pending; waiting map[string]chan bool`。Approve：加锁生成 id、建 `make(chan bool, 1)` 放进 map，解锁后发请求、等答案，两处都和 `ctx.Done()` select。
3. Resolve：加锁取出并删除 map 项，解锁后再往 channel 发送；找不到返回 false。Approve 返回前也要从 map 删掉自己（defer 一个具名方法 forget）。
