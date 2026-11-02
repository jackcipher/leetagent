# 提示

1. 去代码块：`TrimSpace` → 若以 "```" 开头，去掉第一行（```json），再去掉结尾的 "```"。
2. `json.Unmarshal([]byte("[]"), &facts)` 得到的是空切片还是 nil？试一下；必要时 `if facts == nil { facts = []Fact{} }`。
3. Render 复用 All()（已排序），截断到 max，用 strings.Builder 拼接。
