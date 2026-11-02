# 提示

1. 注册表用 `map[string]entry`（entry 存 common.Tool 和 Handler），Register 里懒初始化 map。
2. Serve：`bufio.Scanner` 逐行；`json.Unmarshal` 失败写 parse error（`ID: json.RawMessage("null")`）；`len(req.ID)==0` 是通知，continue；否则 `enc.Encode(s.dispatch(ctx, req))`（Encoder 自动加换行）。
3. dispatch 用 switch method。tools/call：解析 `{name, arguments}`，找不到工具返回 RPCError；Handler 的错误要放进 result 的 isError，而不是 RPCError。
