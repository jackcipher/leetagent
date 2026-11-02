# 提示

1. ExtractCitations：`regexp.MustCompile(`+"`"+`\[(\d+)\]`+"`"+`)` + `FindAllStringSubmatch` + `seen map[int]bool` 去重。
2. Run 的主循环和阶段 4 ex01 一样；把执行单个工具写成具名方法 `execute(ctx, call) common.Message`，出错时 Content = "error: "+err.Error()。
3. 拿到最终回答后再写个 `verify(ctx, &rep)`：对每个引用 `Tools.Call(ctx, "get_post", fmt.Sprintf(`+"`"+`{"tid":%d}`+"`"+`, tid))`，err==nil 放 Verified，否则 Unverified。
