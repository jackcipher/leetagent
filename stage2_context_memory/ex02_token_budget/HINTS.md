# 提示

1. `unicode.IsSpace(r)` 判断空白；向上取整 `(n+3)/4`。
2. 先把 system 拆出来（head），剩下的 body 找出每轮起点：`Role == "user"` 的下标（以及下标 0，防止 body 以 assistant 开头）。
3. 依次尝试从第 k 轮开始保留：`TotalTokens(head)+TotalTokens(body[starts[k]:]) <= budget` 就停；最多试到倒数第二轮，最后一轮无条件保留。用 `make` + `append` 构造新切片，不要原地改。
