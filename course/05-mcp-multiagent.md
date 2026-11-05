# 第 5 课：MCP 与多 Agent

[课程目录](README.md) · [完整阶段讲义](../stage5_mcp_multiagent/README.md)

## 学习目标

- 理解 MCP 的 JSON-RPC 消息与能力协商。
- 实现 MCP Server、Client 和工具调用。
- 用“委派即工具”组织子 Agent，并实现任务路由。

## 课程主线

MCP 解决工具如何被不同客户端复用，多 Agent 解决复杂任务如何拆给不同角色。本课不依赖黑盒框架，直接从协议消息和调用边界开始实现。

## 配套练习

| # | 练习 | 重点 |
|---|---|---|
| 1 | [MCP Server](../stage5_mcp_multiagent/ex01_mcp_server/) | JSON-RPC、stdio、工具暴露 |
| 2 | [MCP Client](../stage5_mcp_multiagent/ex02_mcp_client/) | 握手、工具发现、远程调用 |
| 3 | [子 Agent](../stage5_mcp_multiagent/ex03_subagent/) | 委派、隔离上下文与工具集 |
| 4 | [任务路由](../stage5_mcp_multiagent/ex04_router/) | 规则优先与 LLM 兜底 |

## 完成标准

能够解释 MCP Server 和 Client 的职责，并实现一个可以路由及委派任务的多 Agent 结构。
