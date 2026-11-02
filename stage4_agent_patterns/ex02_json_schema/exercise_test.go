package ex02_json_schema

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"leetagent/common"
)

const reportSchema = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["title", "score", "tags"],
  "properties": {
    "title": {"type": "string", "minLength": 2, "maxLength": 10},
    "score": {"type": "integer", "minimum": 1, "maximum": 5},
    "tags":  {"type": "array", "items": {"type": "string", "enum": ["go", "rust", "ai"]}},
    "draft": {"type": "boolean"},
    "ratio": {"type": "number"}
  }
}`

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		data string
		want []string
	}{
		{"ok", `{"title":"Go 调度","score":4,"tags":["go","ai"],"ratio":0.5}`, nil},
		{"integer as 3.0", `{"title":"标题","score":3.0,"tags":[]}`, nil},
		{"missing", `{"title":"标题"}`, []string{"$.score: required", "$.tags: required"}},
		{"types", `{"title":5,"score":2.5,"tags":"go","draft":"yes"}`, []string{"$.draft: expected boolean", "$.score: expected integer", "$.tags: expected array", "$.title: expected string"}},
		{"ranges", `{"title":"一","score":9,"tags":["java","go"]}`, []string{"$.score: greater than maximum 5", "$.tags[0]: not in enum", "$.title: shorter than 2"}},
		{"extra", `{"title":"标题","score":1,"tags":[],"x":1}`, []string{"$.x: additional property not allowed"}},
		{"root type", `[1]`, []string{"$: expected object"}},
		{"bad json", `{"title":`, []string{"$: invalid json"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Validate(json.RawMessage(reportSchema), json.RawMessage(tc.data))
			if len(got) == 0 && len(tc.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got  %q\nwant %q", got, tc.want)
			}
		})
	}
	if Validate(json.RawMessage(`{"type":"string","maxLength":3}`), json.RawMessage(`"四个汉字"`))[0] != "$: longer than 3" {
		t.Fatal("maxLength should count runes")
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct{ in, want string }{
		{"```json\n{\"a\":1}\n```", `{"a":1}`},
		{"好的，结果是 {\"a\":{\"b\":2}} 希望有帮助", `{"a":{"b":2}}`},
		{"  [1,2] ", `[1,2]`},
	}
	for _, tc := range tests {
		if got := ExtractJSON(tc.in); got != tc.want {
			t.Errorf("ExtractJSON(%q) = %q", tc.in, got)
		}
	}
}

func TestAskStructured(t *testing.T) {
	fake := common.NewScripted(
		common.AssistantText(`{"title":"Go","score":7,"tags":["go"]}`),
		common.AssistantText("```json\n{\"title\":\"Go 帖子\",\"score\":5,\"tags\":[\"go\"]}\n```"),
	)
	msgs := []common.Message{{Role: "user", Content: "给这个帖子打分"}}
	out, err := AskStructured(context.Background(), fake, "glm", msgs, json.RawMessage(reportSchema), 2)
	if err != nil || string(out) != `{"title":"Go 帖子","score":5,"tags":["go"]}` {
		t.Fatalf("out=%s err=%v", out, err)
	}
	second := fake.Calls[1].Msgs
	if len(second) != 3 || second[1].Role != "assistant" || !strings.Contains(second[2].Content, "$.score: greater than maximum 5") {
		t.Fatalf("retry feedback wrong: %+v", second)
	}
	if fake.Calls[0].Tools != nil {
		t.Fatal("no tools for structured output")
	}

	bad := common.NewScripted(common.AssistantText("不知道"), common.AssistantText("还是不知道"))
	if _, err := AskStructured(context.Background(), bad, "m", msgs, json.RawMessage(reportSchema), 1); !errors.Is(err, ErrInvalidAfterRetries) || bad.CallCount() != 2 {
		t.Fatalf("err=%v calls=%d", err, bad.CallCount())
	}
}
