package ex02_model_whitelist

import (
	"errors"
	"reflect"
	"testing"
)

type fakeEnv map[string]string

func (f fakeEnv) Get(k string) string { return f[k] }

const (
	qwen = "mlx-community/Qwen2.5-0.5B-Instruct-4bit"
	glm  = "mlx-community/GLM-4.7-Flash-4bit"
)

func TestLoadModelConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     fakeEnv
		want    ModelConfig
		wantErr error
	}{
		{"default first", fakeEnv{"AGENT_MODELS": " " + qwen + " , ," + glm + "," + qwen}, ModelConfig{Default: qwen, Allowed: []string{qwen, glm}}, nil},
		{"explicit default", fakeEnv{"AGENT_MODELS": qwen + "," + glm, "AGENT_DEFAULT_MODEL": glm}, ModelConfig{Default: glm, Allowed: []string{qwen, glm}}, nil},
		{"empty", fakeEnv{}, ModelConfig{}, ErrNoModels},
		{"bad default", fakeEnv{"AGENT_MODELS": qwen, "AGENT_DEFAULT_MODEL": "gpt-x"}, ModelConfig{}, ErrBadDefault},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LoadModelConfig(tc.env)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %+v, %v; want %+v", got, err, tc.want)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	cfg := ModelConfig{Default: qwen, Allowed: []string{qwen, glm}}
	tests := []struct {
		in, want string
		wantErr  error
	}{
		{"", qwen, nil},
		{glm, glm, nil},
		{"  ", qwen, nil},
		{"gpt-4o", "", ErrModelNotAllowed},
	}
	for _, tc := range tests {
		got, err := cfg.Resolve(tc.in)
		if tc.wantErr != nil {
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Resolve(%q) err = %v", tc.in, err)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("Resolve(%q) = %q, %v; want %q", tc.in, got, err, tc.want)
		}
	}
}
