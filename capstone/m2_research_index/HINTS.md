# 提示

1. Tokenize 维护两个缓冲 `word []rune`、`han []rune`；遇到不同类别字符就 flush 另一个。写个具名函数 `flushHan(run, out)` 生成 bigram。
2. Embed：`h := fnv.New32a(); h.Write([]byte(tok)); v[h.Sum32()%uint32(Dim)]++`，最后除以 sqrt(平方和)。
3. Search：`best := map[int]Hit{}` 保留每个 tid 最高分，再转切片用 sort.Interface 排序。
