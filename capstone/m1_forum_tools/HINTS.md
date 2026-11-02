# 提示

1. 用一个 `args struct{Limit int; TID int; Query string}` 一次解析所有工具的参数，再 switch name。
2. 排序写两个 sort.Interface 类型：byHot（replies 降序、tid 升序）和 byTID。
3. 结果统一 `json.Marshal` 成字符串返回；空结果用 `list := []Summary{}`。
