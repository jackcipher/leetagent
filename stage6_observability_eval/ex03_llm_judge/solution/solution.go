// Package solution 是 ex03_llm_judge 的参考答案。
package solution

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"

	"leetagent/common"
)

var ErrBadVerdict = errors.New("judge output is not a valid verdict")

const JudgeSystem = "你是严格的评审。根据问题、参考答案和待评回答，给待评回答打 1~5 分（5 最好）。" +
	"只输出 JSON：{\"score\": <1-5 的整数>, \"reason\": \"<一句话理由>\"}"

type Verdict struct {
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

type rawVerdict struct {
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

func ParseVerdict(output string) (Verdict, error) {
	i, j := strings.Index(output, "{"), strings.LastIndex(output, "}")
	if i < 0 || j <= i {
		return Verdict{}, ErrBadVerdict
	}
	var rv rawVerdict
	if err := json.Unmarshal([]byte(output[i:j+1]), &rv); err != nil {
		return Verdict{}, ErrBadVerdict
	}
	if rv.Score != math.Trunc(rv.Score) || rv.Score < 1 || rv.Score > 5 {
		return Verdict{}, ErrBadVerdict
	}
	return Verdict{Score: int(rv.Score), Reason: rv.Reason}, nil
}

func JudgeUserPrompt(question, reference, answer string) string {
	return "【问题】\n" + question + "\n\n【参考答案】\n" + reference + "\n\n【待评回答】\n" + answer
}

func Judge(ctx context.Context, c common.Chatter, model, question, reference, answer string) (Verdict, error) {
	msgs := []common.Message{
		{Role: "system", Content: JudgeSystem},
		{Role: "user", Content: JudgeUserPrompt(question, reference, answer)},
	}
	for attempt := 0; attempt < 2; attempt++ {
		reply, err := c.Complete(ctx, model, msgs, nil)
		if err != nil {
			return Verdict{}, err
		}
		if v, err := ParseVerdict(reply.Content); err == nil {
			return v, nil
		}
	}
	return Verdict{}, ErrBadVerdict
}

func Agreement(a, b []int, tolerance int) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	ok := 0
	for i := range a {
		d := a[i] - b[i]
		if d < 0 {
			d = -d
		}
		if d <= tolerance {
			ok++
		}
	}
	return float64(ok) / float64(len(a))
}
