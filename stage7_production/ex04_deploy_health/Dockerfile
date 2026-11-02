# TODO: 把这个"能跑但不合格"的 Dockerfile 改成生产可用：
#   1. 多阶段构建（golang 镜像编译，小镜像运行，例如 gcr.io/distroless/static 或 alpine）
#   2. 加 HEALTHCHECK（探测 /healthz）
#   3. 用非 root 用户运行
FROM golang:1.27
WORKDIR /app
COPY . .
RUN go build -o /app/server ./server
EXPOSE 8080
CMD ["/app/server"]
