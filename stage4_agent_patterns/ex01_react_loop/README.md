# ex01 ReAct 循环 与 Plan-and-Execute（Medium）

## 题目
你在 你的 Agent 项目 里已经写过 `runOneRound`。这里把它抽象成可测试的 `Agent.Run`，再在上面实现"先规划后执行"。

## ReAct（Reason + Act）
```
loop:
  reply = LLM(history, tools)
  history += reply
  if no tool_calls: return reply.content
  for call in reply.tool_calls:
      history += tool(call)          # 出错也回填 "error: ..."
```
## Plan-and-Execute
```
plan = LLM(PlanPrompt, goal)        # 不带工具
for i, step in plan:
    answer_i = ReAct(system="第 i 步: step", user=goal)
```

## 要求
见 exercise.go。重点：
- 工具错误**回填给模型**而不是中断循环——模型常常能换参数重试或向用户解释。
- `MaxRounds` 防止死循环（小模型很容易反复调同一个工具）。
- `TestRunOverMockHTTP` 走真实 SSE 流：参数被切成 2 字符一片，靠 common.Client 拼接。

## 验收
`go test ./stage4_agent_patterns/ex01_react_loop/`

## 对比思考
| | ReAct | Plan-and-Execute |
|---|---|---|
| 调用次数 | 少 | 多（规划 + 每步） |
| 长任务稳定性 | 易跑偏 | 好，步骤可见可审计 |
| 适合 | 简单问答、1~3 次工具 | 调研报告类（毕业项目） |
用 本地小模型 实际跑一下两种方式，记录差异，写进你的阶段总结。
