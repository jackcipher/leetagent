package ex02_chunking

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func texts(cs []Chunk) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.Text)
	}
	return out
}

func TestSplitFixed(t *testing.T) {
	tests := []struct {
		name          string
		in            string
		size, overlap int
		want          []string
		starts        []int
	}{
		{"no overlap", "一二三四五六七", 3, 0, []string{"一二三", "四五六", "七"}, []int{0, 3, 6}},
		{"overlap", "一二三四五六七", 4, 2, []string{"一二三四", "三四五六", "五六七"}, []int{0, 2, 4}},
		{"exact end", "abcdef", 4, 2, []string{"abcd", "cdef"}, []int{0, 2}},
		{"short", "ab", 5, 1, []string{"ab"}, []int{0}},
		{"empty", "", 5, 1, nil, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SplitFixed(tc.in, tc.size, tc.overlap)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(texts(got), tc.want) {
				t.Fatalf("got %q want %q", texts(got), tc.want)
			}
			for i, c := range got {
				if c.Index != i || c.Start != tc.starts[i] {
					t.Fatalf("chunk %d: index=%d start=%d", i, c.Index, c.Start)
				}
			}
		})
	}
	for _, p := range [][2]int{{0, 0}, {3, 3}, {3, -1}} {
		if _, err := SplitFixed("x", p[0], p[1]); !errors.Is(err, ErrBadParams) {
			t.Errorf("params %v: err = %v", p, err)
		}
	}
}

func TestSplitParagraphs(t *testing.T) {
	text := "第一段。\n\n第二段。\n\n\n  \n\n" + strings.Repeat("长", 12) + "\n\n尾段"
	got, err := SplitParagraphs(text, 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"第一段。\n第二段。", "长长长长长长长长长长", "长长长长", "尾段"}
	if !reflect.DeepEqual(texts(got), want) {
		t.Fatalf("got %q", texts(got))
	}
	for i, c := range got {
		if c.Index != i || c.Start != -1 {
			t.Fatalf("chunk %d meta %+v", i, c)
		}
	}
}
