# 第 1 课：Go 工程基础 × LLM 调用

[课程目录](README.md) · [完整阶段讲义](../stage1_foundation/README.md)

## 学习目标

- 理解 messages、tools 和 tool_calls 的请求链路。
- 用接口和注册表组织可扩展工具。
- 正确处理超时、取消、重试、日志和上下文长度。

## 课程主线

一次 Agent 请求并不神秘：组织 messages，调用 LLM，识别 tool_calls，执行工具，再把结果作为 tool 消息发回模型。本课重点是把这条链路写得稳定、可测试、可观测。

## 配套练习

| # | 练习 | 重点 |
|---|---|---|
| 1 | [工具 Registry](../stage1_foundation/ex01_registry/) | interface、map、稳定排序 |
| 2 | [模型白名单](../stage1_foundation/ex02_model_whitelist/) | 配置校验与默认值 |
| 3 | [重试与退避](../stage1_foundation/ex03_retry_backoff/) | context、错误分类、指数退避 |
| 4 | [工具结果截断](../stage1_foundation/ex04_truncate/) | 上下文预算与 UTF-8 安全 |
| 5 | [结构化日志](../stage1_foundation/ex05_slog/) | slog、敏感信息保护 |

## 完成标准

能够解释一次工具调用的完整消息序列，并写出支持取消、有限重试和结构化日志的调用层。
