# 提示

1. 分词：用 strings.Builder 累积 ASCII 字母数字（`r < 128 && (unicode.IsLetter(r) || unicode.IsDigit(r))`），遇到其它字符先把累积的词输出；汉字用 `unicode.Is(unicode.Han, r)` 判断后单独输出。
2. NewBM25 预计算：每篇的 tf map 和长度、全局 df（每篇里出现过就 +1，不是出现次数）、avgLen。
3. 排序沿用 ex03 的 `byScore` sort.Interface 写法；RRF 先累加到 `map[string]float64`，再转切片排序。
