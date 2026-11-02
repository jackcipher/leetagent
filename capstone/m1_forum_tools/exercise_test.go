package m1_forum_tools

import (
	"context"
	"encoding/json"
	"testing"

	"leetagent/common"
)

var _ common.ToolRunner = (*Forum)(nil)

func load(t *testing.T) *Forum {
	t.Helper()
	f, err := LoadForum("testdata/forum.json")
	if err != nil || f == nil || len(f.Posts) != 5 {
		t.Fatalf("LoadForum: %v %v", f, err)
	}
	return f
}

func TestDefs(t *testing.T) {
	f := load(t)
	defs := f.Defs()
	want := []string{"get_list", "get_post", "search_posts"}
	if len(defs) != 3 {
		t.Fatalf("defs = %d", len(defs))
	}
	for i, d := range defs {
		if d.Type != "function" || d.Function.Name != want[i] || d.Function.Description == "" || !json.Valid(d.Function.Parameters) {
			t.Errorf("def %d = %+v", i, d)
		}
	}
}

func TestCall(t *testing.T) {
	f := load(t)
	tests := []struct {
		name, tool, args, want, wantErr string
	}{
		{"list top2", "get_list", `{"limit":2}`, `[{"tid":103,"title":"本地部署 LLM 的显存估算","replies":30},{"tid":101,"title":"Go 语言 goroutine 泄漏排查","replies":12}]`, ""},
		{"list default", "get_list", `{}`, "", ""},
		{"post", "get_post", `{"tid":105}`, `{"tid":105,"title":"goroutine 和线程的区别","author":"dave","replies":3,"content":"goroutine 初始栈只有几 KB，由 Go 调度器在少量系统线程上调度。"}`, ""},
		{"post missing", "get_post", `{"tid":999}`, "", "post 999 not found"},
		{"search case", "search_posts", `{"query":"GOROUTINE"}`, `[{"tid":101,"title":"Go 语言 goroutine 泄漏排查","replies":12},{"tid":105,"title":"goroutine 和线程的区别","replies":3}]`, ""},
		{"search empty result", "search_posts", `{"query":"rust"}`, `[]`, ""},
		{"search no query", "search_posts", `{}`, "", "query is required"},
		{"unknown", "delete_post", `{}`, "", "unknown tool: delete_post"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := f.Call(context.Background(), tc.tool, tc.args)
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("err = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tc.name == "list default" {
				var l []Summary
				if json.Unmarshal([]byte(got), &l) != nil || len(l) != 5 {
					t.Fatalf("got %s", got)
				}
				return
			}
			if got != tc.want {
				t.Fatalf("got  %s\nwant %s", got, tc.want)
			}
		})
	}
	if _, err := f.Call(context.Background(), "get_post", `{tid:1}`); err == nil {
		t.Fatal("bad JSON must error")
	}
}
