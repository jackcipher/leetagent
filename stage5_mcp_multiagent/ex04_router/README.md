# ex04 意图路由：关键词 + LLM 分类兜底（Easy）

## 题目
不是每个请求都需要大 Agent。先**便宜地**判断意图，再分发：

```mermaid
flowchart LR
  T[用户输入] --> K{关键词命中?}
  K -- 是 --> R1[路由 keyword]
  K -- 否 --> L[LLM 分类]
  L --> P{输出可解析?}
  P -- 是 --> R2[路由 llm]
  P -- 否/出错 --> D[Default]
```

## 要求
见 exercise.go。来源字段（keyword/llm/default）方便你统计每种路径的比例。

## 小模型输出的真实样子
你会看到 `Forum。`、`` `code` ``、`类别：forum`、`我认为是 forum`……
所以解析要宽容，但**有歧义时要保守**（回落到 Default）。

## 验收
`go test ./stage5_mcp_multiagent/ex04_router/`
