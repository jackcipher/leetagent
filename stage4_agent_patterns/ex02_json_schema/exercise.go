//go:build !solution

package ex02_json_schema

import (
	"context"
	"encoding/json"
	"errors"

	"leetagent/common"
)

var ErrInvalidAfterRetries = errors.New("structured output still invalid after retries")

// Validate 用 JSON Schema 的一个子集校验 data，返回所有错误（为空表示通过），每条形如 "<path>: <原因>"。
// path 从 "$" 开始，对象字段 "$.a.b"，数组元素 "$.items[0]"。
// 需要支持的关键字：
//
//	type: "object" | "array" | "string" | "integer" | "number" | "boolean" | "null"
//	      （integer 要求是整数值，如 3 或 3.0；number 接受任意数字）
//	properties, required, additionalProperties: false（出现未声明字段报错）
//	items（数组每个元素按 items 校验）
//	enum（值必须等于其中之一，按 JSON 值比较）
//	minimum / maximum（数字）, minLength / maxLength（字符串，按 rune）
//
// 类型不匹配时只报类型错误，不再检查该节点的其它关键字。错误按 path 字典序排序后返回。
// schema 或 data 不是合法 JSON 时返回单条错误 "$: invalid json"。
func Validate(schema, data json.RawMessage) []string {
	// TODO
	return nil
}

// ExtractJSON 从模型输出里取出 JSON：去掉 ``` 代码块；若仍不是以 { 或 [ 开头，取第一个 '{' 到最后一个 '}' 之间的内容。
func ExtractJSON(output string) string {
	// TODO
	return ""
}

// AskStructured 要求模型按 schema 输出 JSON，最多尝试 1+retries 次：
//
//	第一次发送 msgs；之后每次在历史后追加模型上一次的回复（assistant）和
//	{role:"user", content: "你的输出不符合要求：\n" + 错误逐行 + "\n请只输出修正后的 JSON。"}。
//	通过校验返回 ExtractJSON 后的 JSON；用尽仍失败返回 ErrInvalidAfterRetries。
func AskStructured(ctx context.Context, c common.Chatter, model string, msgs []common.Message, schema json.RawMessage, retries int) (json.RawMessage, error) {
	// TODO
	return nil, nil
}
