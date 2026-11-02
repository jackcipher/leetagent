# 提示

1. RenderTranscript 用 switch m.Role；assistant 若 `Content != ""` 才输出内容行（但没有工具调用时即使内容为空也输出），然后对每个 ToolCall 输出箭头行。
2. 找最近 keepTurns 轮的起点：从后往前数 user，数到第 keepTurns 个就是起点下标。
3. 组装结果用三次 append：head、摘要 system、recent。别忘了 `strings.TrimSpace(resp.Content)`。
