# 提示

1. 结构：`dim int` + `map[string]item`，item 存单位化后的向量和原文。
2. 排序不用闭包：定义 `type byScore []Hit`，实现 Len/Swap/Less 三个方法，然后 `sort.Sort(byScore(hits))`。
3. Less：`if a.Score != b.Score { return a.Score > b.Score }; return a.ID < b.ID`。
