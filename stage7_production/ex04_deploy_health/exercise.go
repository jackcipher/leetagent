//go:build !solution

package ex04_deploy_health

import (
	"context"
	"net/http"
	"time"
)

// DockerfilePath 是你要完善的 Dockerfile（相对本目录）。测试会用 CheckDockerfile 检查它。
const DockerfilePath = "Dockerfile"

// Checker 是一个依赖检查（数据库、上游 LLM……）。
type Checker interface {
	Check(ctx context.Context) error
}

// Health 提供 /healthz 和 /readyz。
type Health struct {
	Checks  map[string]Checker
	Timeout time.Duration // 每个检查的超时；0 表示 1 秒
}

// Healthz（存活探针）：进程活着就返回 200，正文 "ok"。不检查依赖——依赖挂了重启进程也没用。
func (h *Health) Healthz(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// Readyz（就绪探针）：按名字排序依次运行所有检查（每个带 Timeout 超时的 ctx），
// 响应 Content-Type: application/json，正文 {"status":"ready"|"not ready","checks":{"<name>":"ok"|"<错误信息>"}}；
// 全部成功 200，否则 503。
func (h *Health) Readyz(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// Env 抽象环境变量。
type Env interface {
	Get(key string) string
}

// Config 是服务配置。
type Config struct {
	Port       int
	LLMBaseURL string
	LLMTimeout time.Duration
	LogLevel   string
}

// LoadConfig 读取并校验配置，遇到第一个错误即返回：
//
//	PORT         默认 8080；必须是 1~65535 的整数，否则 error "PORT: invalid"
//	LLM_BASE_URL 必填（否则 "LLM_BASE_URL: required"）；必须以 http:// 或 https:// 开头（否则 "LLM_BASE_URL: must start with http:// or https://"）；去掉末尾的 "/"
//	LLM_TIMEOUT  默认 60s；用 time.ParseDuration 解析，必须 > 0，否则 "LLM_TIMEOUT: invalid"
//	LOG_LEVEL    默认 "info"；转小写后必须是 debug/info/warn/error 之一，否则 "LOG_LEVEL: invalid"
func LoadConfig(env Env) (Config, error) {
	// TODO
	return Config{}, nil
}

// CheckDockerfile 检查 Dockerfile 内容是否符合生产要求，返回问题列表（按以下顺序）：
//
//	"multi-stage: need at least 2 FROM"       FROM 指令少于 2 条
//	"healthcheck: missing HEALTHCHECK"        没有 HEALTHCHECK 指令
//	"user: must not run as root"              没有 USER 指令，或最后一条 USER 是 root 或 0
//
// 指令名大小写不敏感；忽略空行和 # 注释行；指令判断看每行第一个单词。
func CheckDockerfile(content string) []string {
	// TODO
	return nil
}
