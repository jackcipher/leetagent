# 第 6 课：可观测性与评测

[课程目录](README.md) · [完整阶段讲义](../stage6_observability_eval/README.md)

## 学习目标

- 用 Trace 还原一次 Agent 请求的调用树。
- 建立可重复运行的评测集与评分器。
- 统计 TTFT、P95、吞吐和 Token 成本。

## 课程主线

没有 Trace，很难知道 Agent 慢在哪里；没有评测集，就无法判断一次修改究竟变好还是变坏。本课把“感觉”变成可复现的事件和指标。

## 配套练习

| # | 练习 | 重点 |
|---|---|---|
| 1 | [Agent Trace](../stage6_observability_eval/ex01_trace/) | span 树、属性与错误传播 |
| 2 | [评测跑批](../stage6_observability_eval/ex02_eval_runner/) | JSONL、grader、报告 |
| 3 | [LLM-as-Judge](../stage6_observability_eval/ex03_llm_judge/) | 开放题评分与一致率 |
| 4 | [成本与延迟](../stage6_observability_eval/ex04_cost_latency/) | TTFT、P95、tokens/s |

## 完成标准

能够为 Agent 建立最小可用的 Trace 和回归评测，让优化结果可以用数据验证。
