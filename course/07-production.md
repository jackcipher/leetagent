# 第 7 课：生产化

[课程目录](README.md) · [完整阶段讲义](../stage7_production/README.md)

## 学习目标

- 识别提示注入和不可信工具结果。
- 实现限流、每日额度、SSE 心跳和模型降级。
- 编写健康检查、配置校验和生产容器。

## 课程主线

能跑通的 Agent 不等于能上线的 Agent。本课关注失败、攻击、突发流量、长连接断开和错误配置等生产问题，并用明确的工程边界限制风险。

## 配套练习

| # | 练习 | 重点 |
|---|---|---|
| 1 | [提示注入防护](../stage7_production/ex01_injection_guard/) | 不可信内容、白名单与转义 |
| 2 | [限流与额度](../stage7_production/ex02_rate_limiter/) | 令牌桶、429、用户维度 |
| 3 | [SSE 与降级](../stage7_production/ex03_sse_heartbeat_fallback/) | 心跳、并发写与 fallback |
| 4 | [部署健康](../stage7_production/ex04_deploy_health/) | healthz、readyz、fail fast |

## 完成标准

能够列出 Agent 服务上线前必须具备的安全、流量、可靠性和部署保障，并实现其核心机制。
