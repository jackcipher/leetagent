package ex04_deploy_health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type okCheck struct{}

func (okCheck) Check(ctx context.Context) error { return nil }

type badCheck struct{}

func (badCheck) Check(ctx context.Context) error { return errors.New("connection refused") }

type slowCheck struct{}

func (slowCheck) Check(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return nil
	}
}

func TestHealthz(t *testing.T) {
	h := &Health{Checks: map[string]Checker{"db": badCheck{}}}
	rec := httptest.NewRecorder()
	h.Healthz(rec, httptest.NewRequest("GET", "/healthz", nil))
	if rec.Code != 200 || rec.Body.String() != "ok" {
		t.Fatalf("healthz = %d %q (liveness must not check deps)", rec.Code, rec.Body.String())
	}
}

func TestReadyz(t *testing.T) {
	tests := []struct {
		name   string
		checks map[string]Checker
		code   int
		status string
		want   map[string]string
	}{
		{"all ok", map[string]Checker{"db": okCheck{}, "llm": okCheck{}}, 200, "ready", map[string]string{"db": "ok", "llm": "ok"}},
		{"one bad", map[string]Checker{"db": okCheck{}, "llm": badCheck{}}, 503, "not ready", map[string]string{"db": "ok", "llm": "connection refused"}},
		{"timeout", map[string]Checker{"llm": slowCheck{}}, 503, "not ready", map[string]string{"llm": "context deadline exceeded"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := &Health{Checks: tc.checks, Timeout: 20 * time.Millisecond}
			rec := httptest.NewRecorder()
			start := time.Now()
			h.Readyz(rec, httptest.NewRequest("GET", "/readyz", nil))
			if time.Since(start) > time.Second {
				t.Fatal("checks must respect Timeout")
			}
			var body struct {
				Status string            `json:"status"`
				Checks map[string]string `json:"checks"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body %q: %v", rec.Body.String(), err)
			}
			if rec.Code != tc.code || body.Status != tc.status || !reflect.DeepEqual(body.Checks, tc.want) || rec.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("code=%d body=%+v ct=%q", rec.Code, body, rec.Header().Get("Content-Type"))
			}
		})
	}
}

type mapEnv map[string]string

func (m mapEnv) Get(k string) string { return m[k] }

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig(mapEnv{"LLM_BASE_URL": "http://localhost:8000/v1/"})
	want := Config{Port: 8080, LLMBaseURL: "http://localhost:8000/v1", LLMTimeout: 60 * time.Second, LogLevel: "info"}
	if err != nil || cfg != want {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
	cfg, err = LoadConfig(mapEnv{"PORT": "9000", "LLM_BASE_URL": "https://x", "LLM_TIMEOUT": "90s", "LOG_LEVEL": "DEBUG"})
	if err != nil || cfg.Port != 9000 || cfg.LLMTimeout != 90*time.Second || cfg.LogLevel != "debug" {
		t.Fatalf("cfg=%+v err=%v", cfg, err)
	}
	bad := []struct {
		env  mapEnv
		want string
	}{
		{mapEnv{"PORT": "70000", "LLM_BASE_URL": "http://x"}, "PORT: invalid"},
		{mapEnv{"PORT": "abc"}, "PORT: invalid"},
		{mapEnv{}, "LLM_BASE_URL: required"},
		{mapEnv{"LLM_BASE_URL": "localhost:8000"}, "LLM_BASE_URL: must start with http:// or https://"},
		{mapEnv{"LLM_BASE_URL": "http://x", "LLM_TIMEOUT": "-1s"}, "LLM_TIMEOUT: invalid"},
		{mapEnv{"LLM_BASE_URL": "http://x", "LOG_LEVEL": "verbose"}, "LOG_LEVEL: invalid"},
	}
	for _, tc := range bad {
		if _, err := LoadConfig(tc.env); err == nil || err.Error() != tc.want {
			t.Errorf("%v: err = %v, want %q", tc.env, err, tc.want)
		}
	}
}

func TestCheckDockerfile(t *testing.T) {
	tests := []struct {
		name, content string
		want          []string
	}{
		{"bad", "FROM golang:1.27\n# USER app\nRUN go build\n", []string{"multi-stage: need at least 2 FROM", "healthcheck: missing HEALTHCHECK", "user: must not run as root"}},
		{"root at end", "from a as b\nfrom c\nuser app\nUSER root\nhealthcheck CMD true\n", []string{"user: must not run as root"}},
		{"uid 0", "FROM a\nFROM b\nUSER 0:0\nHEALTHCHECK CMD true\n", []string{"user: must not run as root"}},
		{"good", "FROM a AS build\n\nFROM b\nUSER 10001\nHEALTHCHECK CMD wget -qO- localhost/healthz\n", nil},
	}
	for _, tc := range tests {
		got := CheckDockerfile(tc.content)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %q", tc.name, got)
		}
	}
}

// 这个测试检查的是你写的 Dockerfile 本身。
func TestYourDockerfile(t *testing.T) {
	b, err := os.ReadFile(DockerfilePath)
	if err != nil {
		t.Fatal(err)
	}
	if issues := CheckDockerfile(string(b)); len(issues) != 0 {
		t.Fatalf("%s has issues: %q", DockerfilePath, issues)
	}
	// 防止 CheckDockerfile 是空实现时这个测试白白通过
	if len(CheckDockerfile("FROM x")) == 0 {
		t.Fatal("CheckDockerfile must report problems for a bad Dockerfile")
	}
}
