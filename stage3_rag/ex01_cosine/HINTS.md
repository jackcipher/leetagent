# 提示

1. `math.Sqrt(sum(x*x))`。Cosine 先调 Dot（顺便检查维度），再检查两个范数是否为 0。
2. Normalize 用 `make([]float64, len(v))` 新建切片，不要原地除。
3. 解析时先建 `out := make([][]float64, len(data))`，遍历 data 把 embedding 放到 `out[index]`；index 越界或该位置已有值就是坏响应。
