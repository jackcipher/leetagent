package ex02_eval_runner

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoadCases(t *testing.T) {
	f, err := os.Open("testdata/cases.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cases, err := LoadCases(f)
	if err != nil || len(cases) != 3 {
		t.Fatalf("cases=%d err=%v", len(cases), err)
	}
	if cases[1].ID != "post-7" || cases[1].ExpectTool != "get_post" || cases[1].MustNotContain[0] != "我无法" {
		t.Fatalf("case = %+v", cases[1])
	}
	_, err = LoadCases(strings.NewReader("{\"id\":\"a\",\"input\":\"x\"}\n\n{bad"))
	if err == nil || !strings.HasPrefix(err.Error(), "line 3: ") {
		t.Fatalf("err = %v", err)
	}
	_, err = LoadCases(strings.NewReader(`{"id":"a"}`))
	if err == nil || err.Error() != "line 1: id and input are required" {
		t.Fatalf("err = %v", err)
	}
}

func TestGrade(t *testing.T) {
	c := Case{MustContain: []string{"GMP", "调度"}, MustNotContain: []string{"I don't know"}, ExpectTool: "get_post"}
	tests := []struct {
		out  Output
		want []string
	}{
		{Output{"gmp 模型与调度器", []string{"get_list", "get_post"}}, nil},
		{Output{"I DON'T KNOW", nil}, []string{"missing: GMP", "missing: 调度", "forbidden: I don't know", "tool not used: get_post"}},
	}
	for _, tc := range tests {
		got := Grade(c, tc.out)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got %q", got)
		}
	}
}

type fakeAgent struct{}

func (fakeAgent) Answer(ctx context.Context, input string) (Output, error) {
	switch {
	case strings.Contains(input, "第一页"):
		return Output{"有 GMP、Rust 两个帖子", []string{"get_list"}}, nil
	case strings.Contains(input, "tid 7"):
		return Output{}, errors.New("upstream 502")
	}
	return Output{Text: "抱歉，我不能透露。"}, nil
}

func TestRunAndReport(t *testing.T) {
	f, _ := os.Open("testdata/cases.jsonl")
	defer f.Close()
	cases, _ := LoadCases(f)
	rep := Run(context.Background(), fakeAgent{}, cases)
	if rep.Total != 3 || rep.Passed != 2 || len(rep.Results) != 3 {
		t.Fatalf("rep = %+v", rep)
	}
	if rep.Results[1].Pass || !reflect.DeepEqual(rep.Results[1].Reasons, []string{"error: upstream 502"}) {
		t.Fatalf("result[1] = %+v", rep.Results[1])
	}
	if r := rep.PassRate(); r < 0.666 || r > 0.667 {
		t.Fatalf("rate = %v", r)
	}
	want := "| id | 结果 | 原因 |\n|---|---|---|\n" +
		"| list-1 | ✅ |  |\n" +
		"| post-7 | ❌ | error: upstream 502 |\n" +
		"| safety-1 | ✅ |  |\n" +
		"\n通过 2/3（66.7%）\n"
	if got := rep.Markdown(); got != want {
		t.Fatalf("got:\n%s", got)
	}
	if (Report{}).PassRate() != 0 {
		t.Fatal("empty report rate should be 0")
	}
}
