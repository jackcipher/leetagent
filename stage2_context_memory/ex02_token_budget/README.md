# ex02 Token 估算与预算裁剪（Medium）

## 题目
上下文窗口是 Agent 最稀缺的资源。多轮对话 + 工具结果很快超过模型窗口（GLM-4.7-Flash 本地部署常常只给 8k~32k）。
实现一个不依赖分词器的估算器，和一个"按轮丢最早对话"的裁剪器。

## 估算规则
- CJK 字符（rune >= 0x2E80）每个 1 token
- 其它非空白字符累计，每 4 个 1 token，向上取整
- 空白不计
- 每条消息额外 `MessageOverhead = 4`；工具调用的 name 和 arguments 也要计入

## 裁剪规则（TrimToBudget）
1. 开头的 system 永远保留。
2. 以"轮"为单位丢：user 开始，直到下一个 user 之前。**绝不能把 tool 消息和它对应的 assistant.tool_calls 拆开**——
   否则上游会报 `tool_call_id not found` 之类的 400。
3. 从最早的轮丢，直到 <= budget。
4. 最后一轮必须保留（用户刚问的问题不能丢）。
5. 不修改入参。

## 验收
`go test ./stage2_context_memory/ex02_token_budget/`

## 延伸
精确计数要用模型自己的 tokenizer（tiktoken / HF tokenizers）。上游响应里的 `usage.prompt_tokens` 是最可靠的"事后"数字，可以用来校准估算系数。
