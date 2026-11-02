# 提示

1. 先 `history := make([]common.Message, len(msgs)); copy(history, msgs)`，之后只 append history，入参就不会被改。
2. 把"执行一个工具调用 → 返回 tool 消息"写成方法 `func (a *Agent) runTool(ctx, call) common.Message`，循环体会很清爽。
3. ParsePlan：逐行 TrimSpace，扫描开头的数字（至少一位），后面紧跟 "." "、" ")" 之一，剩余部分 TrimSpace 后非空才算步骤。
