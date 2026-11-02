// Package solution 是 ex04_deploy_health 的参考答案。
package solution

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DockerfilePath 相对于练习目录（测试在练习目录运行）。
const DockerfilePath = "solution/Dockerfile"

type Checker interface {
	Check(ctx context.Context) error
}

type Health struct {
	Checks  map[string]Checker
	Timeout time.Duration
}

func (h *Health) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type readyBody struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

func (h *Health) runCheck(ctx context.Context, c Checker) error {
	timeout := h.Timeout
	if timeout <= 0 {
		timeout = time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Check(ctx)
}

func (h *Health) Readyz(w http.ResponseWriter, r *http.Request) {
	names := make([]string, 0, len(h.Checks))
	for n := range h.Checks {
		names = append(names, n)
	}
	sort.Strings(names)
	body := readyBody{Status: "ready", Checks: map[string]string{}}
	code := http.StatusOK
	for _, n := range names {
		if err := h.runCheck(r.Context(), h.Checks[n]); err != nil {
			body.Checks[n] = err.Error()
			body.Status = "not ready"
			code = http.StatusServiceUnavailable
		} else {
			body.Checks[n] = "ok"
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

type Env interface {
	Get(key string) string
}

type Config struct {
	Port       int
	LLMBaseURL string
	LLMTimeout time.Duration
	LogLevel   string
}

func LoadConfig(env Env) (Config, error) {
	cfg := Config{Port: 8080, LLMTimeout: 60 * time.Second, LogLevel: "info"}
	if v := env.Get("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 || p > 65535 {
			return Config{}, errors.New("PORT: invalid")
		}
		cfg.Port = p
	}
	u := env.Get("LLM_BASE_URL")
	if u == "" {
		return Config{}, errors.New("LLM_BASE_URL: required")
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return Config{}, errors.New("LLM_BASE_URL: must start with http:// or https://")
	}
	cfg.LLMBaseURL = strings.TrimRight(u, "/")
	if v := env.Get("LLM_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			return Config{}, errors.New("LLM_TIMEOUT: invalid")
		}
		cfg.LLMTimeout = d
	}
	if v := env.Get("LOG_LEVEL"); v != "" {
		lv := strings.ToLower(v)
		switch lv {
		case "debug", "info", "warn", "error":
			cfg.LogLevel = lv
		default:
			return Config{}, errors.New("LOG_LEVEL: invalid")
		}
	}
	return cfg, nil
}

func CheckDockerfile(content string) []string {
	froms, health := 0, false
	lastUser := ""
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		switch strings.ToUpper(fields[0]) {
		case "FROM":
			froms++
		case "HEALTHCHECK":
			health = true
		case "USER":
			if len(fields) > 1 {
				lastUser = fields[1]
			}
		}
	}
	var issues []string
	if froms < 2 {
		issues = append(issues, "multi-stage: need at least 2 FROM")
	}
	if !health {
		issues = append(issues, "healthcheck: missing HEALTHCHECK")
	}
	u := strings.SplitN(lastUser, ":", 2)[0]
	if u == "" || u == "root" || u == "0" {
		issues = append(issues, "user: must not run as root")
	}
	return issues
}
