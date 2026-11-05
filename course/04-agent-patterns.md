# 第 4 课：Agent 核心模式

[课程目录](README.md) · [完整阶段讲义](../stage4_agent_patterns/README.md)

## 学习目标

- 实现 ReAct 与 Plan-and-Execute 循环。
- 处理结构化输出、并行工具和最大轮数。
- 为危险操作增加人工审批。

## 课程主线

Agent 的主体通常就是一个受约束的循环。本课把循环中的终止条件、工具执行、并发、结构化输出失败和人工审批逐一拆开，使其能够被单元测试覆盖。

## 配套练习

| # | 练习 | 重点 |
|---|---|---|
| 1 | [ReAct 循环](../stage4_agent_patterns/ex01_react_loop/) | 历史维护、工具回填、轮数限制 |
| 2 | [JSON Schema](../stage4_agent_patterns/ex02_json_schema/) | 提取、校验与纠错重试 |
| 3 | [并行工具](../stage4_agent_patterns/ex03_parallel_tools/) | 并发限制、顺序与取消 |
| 4 | [审批闸门](../stage4_agent_patterns/ex04_approval_gate/) | Human-in-the-loop |
| 5 | [Prompt 模板](../stage4_agent_patterns/ex05_prompt_template/) | 模板管理与工具描述质量 |

## 完成标准

能够写出有终止边界、可并发、可审批、并能处理模型非结构化输出的 Agent 主循环。
