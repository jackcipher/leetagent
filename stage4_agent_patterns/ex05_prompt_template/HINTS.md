# 提示

1. Build：收集非空的 `"## "+Title+"\n"+TrimSpace(Body)`，最后 `strings.Join(parts, "\n\n")`。
2. RenderToolGuide：先复制切片再排序（不能改入参）。按要求不用闭包，定义 `type byName []common.Tool` 实现 sort.Interface。
3. Lint：用 `regexp.MustCompile("^[a-z][a-z0-9_]*$")`；长度用 `utf8.RuneCountInString`；参数解析到 `struct{ Properties map[string]struct{Description string} }`，Unmarshal 失败就是 invalid json；收集缺描述的参数名后 sort.Strings。
