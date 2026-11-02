package ex04_truncate

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name string
		in   string
		max  int
		want string
	}{
		{"short", "你好", 5, "你好"},
		{"exact", "你好世界", 4, "你好世界"},
		{"chinese cut", "今天天气真不错", 3, "今天天" + fmt.Sprintf(TruncMarker, 4)},
		{"ascii", "abcdefghij", 4, "abcd" + fmt.Sprintf(TruncMarker, 6)},
		{"no limit", "abc", 0, "abc"},
		{"emoji", "😀😀😀", 1, "😀" + fmt.Sprintf(TruncMarker, 2)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Truncate(tc.in, tc.max)
			if got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
			if !utf8.ValidString(got) {
				t.Fatal("invalid utf8")
			}
		})
	}
}

func TestTruncateHeadTail(t *testing.T) {
	in := strings.Repeat("甲", 10) + strings.Repeat("乙", 100) + strings.Repeat("丙", 5)
	got := TruncateHeadTail(in, 10, 5)
	want := strings.Repeat("甲", 10) + fmt.Sprintf(MiddleMarker, 100) + strings.Repeat("丙", 5)
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if TruncateHeadTail("短文本", 2, 2) != "短文本" {
		t.Fatal("short text should be unchanged")
	}
}

func TestTruncateToolResult(t *testing.T) {
	limits := map[string]int{"fetch_one": 3}
	if got := TruncateToolResult("fetch_one", "一二三四五", limits, 100); got != "一二三"+fmt.Sprintf(TruncMarker, 2) {
		t.Errorf("fetch_one: %q", got)
	}
	if got := TruncateToolResult("get_list", "一二三四五", limits, 4); got != "一二三四"+fmt.Sprintf(TruncMarker, 1) {
		t.Errorf("default: %q", got)
	}
}
