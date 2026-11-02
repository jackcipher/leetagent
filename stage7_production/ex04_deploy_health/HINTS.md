# 提示

1. Readyz：名字排序；每个检查 `context.WithTimeout(r.Context(), timeout)`（写个具名方法 runCheck 里 defer cancel）；先设 Content-Type 再 WriteHeader，最后 json.Encode。
2. LoadConfig 按 PORT → LLM_BASE_URL → LLM_TIMEOUT → LOG_LEVEL 顺序校验，第一个错误就返回；`strings.TrimRight(u, "/")`。
3. CheckDockerfile：`strings.Fields(line)`，第一个词 ToUpper 后 switch；USER 取第二个词，`USER 0:0` 要取冒号前。Dockerfile 参考：两个 FROM（`golang:1.27-alpine AS build` + `alpine`）、`adduser -D -u 10001 app`、`USER app`、`HEALTHCHECK CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1`。
