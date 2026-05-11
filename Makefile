APP_NAME  := enterprise-agent
GATEWAY   := ./module-gateway
BIN_DIR   := ./tmp
BIN       := $(BIN_DIR)/main
MANIFEST  := ./manifest
DOCKER    := $(MANIFEST)/docker
CONFIG    := $(MANIFEST)/config
RESOURCE  := ./resource

.PHONY: all build run dev clean tidy lint swagger docker docker-up docker-down help

## 默认目标
all: build

## ── 编译 ─────────────────────────────────────────────────────────────────────

build: ## 编译 gateway 服务到 tmp/main
	@mkdir -p $(BIN_DIR)
	cd $(GATEWAY) && go build -o ../$(BIN) .

## ── 运行 ─────────────────────────────────────────────────────────────────────

run: build ## 编译并运行（生产模式）
	APP_ENV=prod $(BIN)

dev: build ## 编译并运行（开发模式）
	APP_ENV=dev $(BIN)

## ── 热重载（Air）────────────────────────────────────────────────────────────

air: ## 使用 Air 启动热重载开发服务
	air -c .air.toml

## ── 依赖管理 ─────────────────────────────────────────────────────────────────

tidy: ## 对所有子模块执行 go mod tidy
	@for mod in module-infra module-config module-agent module-knowledge module-conversation module-gateway; do \
		echo "→ tidy $$mod"; \
		(cd $$mod && go mod tidy); \
	done

## ── 代码生成 ─────────────────────────────────────────────────────────────────

swagger: ## 生成 Swagger 文档（自动安装 swag）
	@command -v swag >/dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	cd $(GATEWAY) && swag init -g main.go -o docs \
		--dir .,../module-agent/handler/agent,../module-knowledge/handler/knowledge,../module-conversation/handler/conversation

wire: ## 重新生成所有 Wire 代码（自动安装 wire）
	@command -v wire >/dev/null 2>&1 || go install github.com/google/wire/cmd/wire@latest
	@for mod in module-agent module-knowledge module-conversation module-core module-gateway; do \
		echo "→ wire $$mod"; \
		(cd $$mod/wire && wire); \
	done

## ── 测试 ─────────────────────────────────────────────────────────────────────

test: ## 运行所有模块的单元测试
	@for mod in module-infra module-config module-agent module-knowledge module-conversation module-gateway; do \
		echo "→ test $$mod"; \
		(cd $$mod && go test ./... -count=1); \
	done

test-cover: ## 生成测试覆盖率报告
	@mkdir -p $(BIN_DIR)
	@for mod in module-agent module-knowledge module-conversation; do \
		echo "→ cover $$mod"; \
		(cd $$mod && go test ./... -coverprofile=../$(BIN_DIR)/$$mod.coverage.out); \
	done

## ── 清理 ─────────────────────────────────────────────────────────────────────

clean: ## 清除编译产物
	@rm -rf $(BIN_DIR)
	@echo "cleaned."

## ── Docker ───────────────────────────────────────────────────────────────────

docker: ## 构建 Docker 镜像
	docker build -f $(DOCKER)/Dockerfile -t $(APP_NAME):latest .

docker-up: ## 启动所有服务（docker compose）
	docker compose -f $(DOCKER)/docker-compose.yaml up -d

docker-down: ## 停止所有服务
	docker compose -f $(DOCKER)/docker-compose.yaml down

docker-logs: ## 查看 gateway 容器日志
	docker compose -f $(DOCKER)/docker-compose.yaml logs -f gateway

## ── 帮助 ─────────────────────────────────────────────────────────────────────

help: ## 打印所有可用 target 说明
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'
