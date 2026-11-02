# 提示

1. KeywordRoute：text 和关键词都 ToLower；每个路由用一个 `map[string]bool` 去重；只有 `n > bestN` 时才更新（这样并列保留靠前的）。
2. 解析：`strings.ToLower(strings.Trim(strings.TrimSpace(out), "`。. \n"))`，先精确匹配，再统计包含了几个路由名。
3. Route 返回两个值；LLM 的 err 不往外抛，直接 Default。
