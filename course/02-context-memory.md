# 第 2 课：上下文与记忆

[课程目录](README.md) · [完整阶段讲义](../stage2_context_memory/README.md)

## 学习目标

- 区分工作记忆、会话记忆和长期记忆。
- 在 Token 预算内保留必要上下文。
- 使用 SQLite 持久化消息，并通过摘要压缩旧对话。

## 课程主线

LLM 没有天然记忆。Agent 的“记得”来自工程系统对消息、摘要和事实的存储与选择。本课关注哪些内容应该保存、何时压缩，以及如何避免历史无限增长。

## 配套练习

| # | 练习 | 重点 |
|---|---|---|
| 1 | [SQLite 会话持久化](../stage2_context_memory/ex01_session_store/) | schema、事务、并发访问 |
| 2 | [Token 预算裁剪](../stage2_context_memory/ex02_token_budget/) | 估算、成组保留消息 |
| 3 | [对话摘要压缩](../stage2_context_memory/ex03_summary_compress/) | 摘要时机与事实保真 |
| 4 | [长期记忆](../stage2_context_memory/ex04_long_term_memory/) | 写入、检索、去重、遗忘 |

## 完成标准

能够设计一个不会无限膨胀、不会拆散工具调用、并能跨会话保留关键事实的上下文系统。
