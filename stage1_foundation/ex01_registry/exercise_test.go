package ex01_registry

import (
	"context"
	"errors"
	"strings"
	"testing"

	"leetagent/common"
)

type echoTool struct{ name string }

func (e echoTool) Def() common.Tool {
	return common.NewFunctionTool(e.name, "echo", `{"type":"object"}`)
}

func (e echoTool) Run(ctx context.Context, args string) (string, error) {
	return e.name + ":" + args, nil
}

func TestRegisterAndDefsSorted(t *testing.T) {
	r := NewRegistry()
	for _, n := range []string{"get_list", "fetch_one", "get_current_time"} {
		if err := r.Register(echoTool{n}); err != nil {
			t.Fatalf("register %s: %v", n, err)
		}
	}
	defs := r.Defs()
	want := []string{"fetch_one", "get_current_time", "get_list"}
	if len(defs) != len(want) {
		t.Fatalf("Defs len = %d, want %d", len(defs), len(want))
	}
	for i, d := range defs {
		if d.Function.Name != want[i] {
			t.Errorf("Defs[%d] = %s, want %s", i, d.Function.Name, want[i])
		}
	}
}

func TestRegisterErrors(t *testing.T) {
	r := NewRegistry()
	if err := r.Register(echoTool{""}); !errors.Is(err, ErrEmptyName) {
		t.Errorf("empty name: got %v", err)
	}
	if err := r.Register(echoTool{"a"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(echoTool{"a"}); !errors.Is(err, ErrDuplicateTool) {
		t.Errorf("duplicate: got %v", err)
	}
}

func TestCall(t *testing.T) {
	r := NewRegistry()
	_ = r.Register(echoTool{"get_list"})
	tests := []struct {
		name, tool, args, want string
		wantErr                error
	}{
		{"known", "get_list", `{"page":1}`, `get_list:{"page":1}`, nil},
		{"unknown", "rm_rf", `{}`, "", ErrUnknownTool},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := r.Call(context.Background(), tc.tool, tc.args)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) || !strings.Contains(err.Error(), tc.tool) {
					t.Fatalf("err = %v, want %v mentioning %s", err, tc.wantErr, tc.tool)
				}
				return
			}
			if err != nil || got != tc.want {
				t.Fatalf("got %q, %v; want %q", got, err, tc.want)
			}
		})
	}
}
