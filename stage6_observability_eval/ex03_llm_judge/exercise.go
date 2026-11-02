//go:build !solution

package ex03_llm_judge

import (
	"context"
	"errors"

	"leetagent/common"
)

var ErrBadVerdict = errors.New("judge output is not a valid verdict")

// JudgeSystem 是裁判的 system prompt。
const JudgeSystem = "你是严格的评审。根据问题、参考答案和待评回答，给待评回答打 1~5 分（5 最好）。" +
	"只输出 JSON：{\"score\": <1-5 的整数>, \"reason\": \"<一句话理由>\"}"

// Verdict 是裁判结论。
type Verdict struct {
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

// ParseVerdict 解析裁判输出：容忍 ``` 代码块和前后废话（取第一个 '{' 到最后一个 '}'）；
// score 必须是 1~5 的整数（JSON 里写成 4.0 也接受，4.5 不接受），否则返回 ErrBadVerdict。
func ParseVerdict(output string) (Verdict, error) {
	// TODO
	return Verdict{}, nil
}

// JudgeUserPrompt 生成 user 消息，格式严格为：
//
//	"【问题】\n<q>\n\n【参考答案】\n<ref>\n\n【待评回答】\n<ans>"
func JudgeUserPrompt(question, reference, answer string) string {
	// TODO
	return ""
}

// Judge 调用模型打分（不带工具）。解析失败时重试一次（再调一次模型，消息相同）；仍失败返回 ErrBadVerdict。
func Judge(ctx context.Context, c common.Chatter, model, question, reference, answer string) (Verdict, error) {
	// TODO
	return Verdict{}, nil
}

// Agreement 计算两组打分的一致率：|a[i]-b[i]| <= tolerance 的比例。长度不同或为空返回 0。
// 用途：拿 LLM 裁判的分数和你自己人工打的分数比，判断裁判是否可信。
func Agreement(a, b []int, tolerance int) float64 {
	// TODO
	return 0
}
