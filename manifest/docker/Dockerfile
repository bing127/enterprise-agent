# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 安装必要工具
RUN apk add --no-cache git

# 拷贝 go workspace 配置和所有模块的 go.mod/go.sum，利用 Docker 缓存层
COPY go.work go.work.sum ./
COPY module-agent/go.mod      module-agent/go.sum      ./module-agent/
COPY module-config/go.mod     module-config/go.sum     ./module-config/
COPY module-conversation/go.mod module-conversation/go.sum ./module-conversation/
COPY module-gateway/go.mod    module-gateway/go.sum    ./module-gateway/
COPY module-infra/go.mod      module-infra/go.sum      ./module-infra/
COPY module-knowledge/go.mod  module-knowledge/go.sum  ./module-knowledge/
COPY module-pkg/go.mod        module-pkg/go.sum        ./module-pkg/

RUN go work sync

# 拷贝所有源码
COPY . .

# 编译 gateway（静态链接，去掉 CGO）
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /build/main \
    ./module-gateway

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

# 安装 CA 证书和时区数据
RUN apk add --no-cache ca-certificates tzdata

# 拷贝二进制
COPY --from=builder /build/main .

# 拷贝配置文件
COPY manifest/config/config.dev.yaml     ./config/config.dev.yaml
COPY manifest/config/config.prod.yaml    ./config/config.prod.yaml

# 运行用户（非 root）
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN mkdir -p /app/resource/logs && chown -R appuser:appgroup /app/resource/logs /app/config
USER appuser

EXPOSE 8080

ENV APP_ENV=prod \
    CONFIG_DIR=/app/config

ENTRYPOINT ["./main"]
