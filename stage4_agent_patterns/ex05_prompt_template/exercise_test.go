package ex05_prompt_template

import (
	"reflect"
	"strings"
	"testing"

	"leetagent/common"
)

func TestBuild(t *testing.T) {
	got := Build([]Section{
		{"角色", "  你是论坛研究助手。\n"},
		{"空节", "   "},
		{"规则", "只用工具返回的数据。\n引用帖子用 [tid]。"},
	})
	want := "## 角色\n你是论坛研究助手。\n\n## 规则\n只用工具返回的数据。\n引用帖子用 [tid]。"
	if got != want {
		t.Fatalf("got:\n%s", got)
	}
	if Build(nil) != "" || Build([]Section{{"x", ""}}) != "" {
		t.Fatal("empty should be empty")
	}
}

func TestRenderToolGuide(t *testing.T) {
	tools := []common.Tool{
		common.NewFunctionTool("get_post", "读取单个帖子", `{}`),
		common.NewFunctionTool("get_list", "列出帖子", `{}`),
	}
	if got := RenderToolGuide(tools); got != "- get_list: 列出帖子\n- get_post: 读取单个帖子" {
		t.Fatalf("got %q", got)
	}
	if tools[0].Function.Name != "get_post" {
		t.Fatal("input must not be reordered")
	}
}

func TestLint(t *testing.T) {
	tests := []struct {
		name string
		tool common.Tool
		want []string
	}{
		{"good", common.NewFunctionTool("get_post", "按 tid 读取帖子正文和回复",
			`{"type":"object","properties":{"tid":{"type":"integer","description":"帖子 ID"}}}`), nil},
		{"bad name short desc", common.NewFunctionTool("GetPost", "读帖子",
			`{"properties":{"tid":{"type":"integer"},"page":{"type":"integer"}}}`),
			[]string{"name: must be snake_case", "description: too short", "parameters.page: missing description", "parameters.tid: missing description"}},
		{"long + invalid", common.NewFunctionTool("x", strings.Repeat("长", 201), `[1]`),
			[]string{"description: too long", "parameters: invalid json"}},
	}
	for _, tc := range tests {
		got := LintToolDescription(tc.tool)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %q", tc.name, got)
		}
	}
}
