package ex04_search_tool

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type fakeEmbedder struct{ got []string }

func (f *fakeEmbedder) Embed(ctx context.Context, text string) ([]float64, error) {
	f.got = append(f.got, text)
	return []float64{1, 0}, nil
}

type fakeIndex struct {
	hits  []Hit
	lastK int
}

func (f *fakeIndex) Search(q []float64, k int) ([]Hit, error) {
	f.lastK = k
	if k < len(f.hits) {
		return f.hits[:k], nil
	}
	return f.hits, nil
}

func TestDef(t *testing.T) {
	s := &SearchPosts{}
	d := s.Def()
	if d.Type != "function" || d.Function.Name != "search_posts" || len([]rune(d.Function.Description)) < 10 {
		t.Fatalf("bad def %+v", d)
	}
	var schema struct {
		Type       string                     `json:"type"`
		Properties map[string]json.RawMessage `json:"properties"`
		Required   []string                   `json:"required"`
	}
	if err := json.Unmarshal(d.Function.Parameters, &schema); err != nil {
		t.Fatalf("parameters not json: %v", err)
	}
	if schema.Type != "object" || schema.Properties["query"] == nil || schema.Properties["top_k"] == nil ||
		len(schema.Required) != 1 || schema.Required[0] != "query" {
		t.Fatalf("bad schema %+v", schema)
	}
}

func TestRun(t *testing.T) {
	idx := &fakeIndex{hits: []Hit{{"t1#0", 0.9, "Go 协程"}, {"t2#3", 0.5, "channel"}, {"t3#1", 0.1, "x"}}}
	emb := &fakeEmbedder{}
	s := &SearchPosts{Embedder: emb, Index: idx, DefaultK: 2, MaxK: 3}
	tests := []struct {
		name   string
		args   string
		wantK  int
		wantN  int
		errIs  error
		anyErr bool
	}{
		{"default k", `{"query":" goroutine "}`, 2, 2, nil, false},
		{"explicit k", `{"query":"go","top_k":1}`, 1, 1, nil, false},
		{"capped", `{"query":"go","top_k":50}`, 3, 3, nil, false},
		{"empty", `{"query":"   "}`, 0, 0, ErrEmptyQuery, true},
		{"bad json", `{"query":`, 0, 0, nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			idx.lastK = 0
			out, err := s.Run(context.Background(), tc.args)
			if tc.anyErr {
				if err == nil || (tc.errIs != nil && !errors.Is(err, tc.errIs)) {
					t.Fatalf("err = %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var hits []Hit
			if err := json.Unmarshal([]byte(out), &hits); err != nil || len(hits) != tc.wantN || idx.lastK != tc.wantK {
				t.Fatalf("out=%s k=%d err=%v", out, idx.lastK, err)
			}
		})
	}
	if emb.got[0] != "goroutine" {
		t.Errorf("query should be trimmed, got %q", emb.got[0])
	}
	empty := &SearchPosts{Embedder: emb, Index: &fakeIndex{}, DefaultK: 3}
	if out, err := empty.Run(context.Background(), `{"query":"x"}`); err != nil || out != "[]" {
		t.Errorf("no hits: %q %v", out, err)
	}
}
