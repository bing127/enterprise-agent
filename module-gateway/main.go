// Package main is the entry point for the Enterprise Agent gateway service.
//
// @title           Enterprise Agent API
// @version         1.0
// @description     企业级智能 Agent 系统 API。HTTP 层基于 CloudWeGo Hertz，AI 对话层基于 CloudWeGo Eino。
//
// @contact.name    Enterprise Agent Team
// @contact.email   support@example.com
//
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 格式：Bearer <token>
package main

import (
	"os"
	"time"

	agentHandler "github.com/bing127/enterprise-agent/module-agent/handler/agent"
	agentService "github.com/bing127/enterprise-agent/module-agent/service"
	conversationHandler "github.com/bing127/enterprise-agent/module-conversation/handler/conversation"
	conversationService "github.com/bing127/enterprise-agent/module-conversation/service"
	coreHandler "github.com/bing127/enterprise-agent/module-core/handler"
	coreRepo "github.com/bing127/enterprise-agent/module-core/repo"
	coreService "github.com/bing127/enterprise-agent/module-core/service"
	"github.com/bing127/enterprise-agent/module-gateway/router"
	knowledgeHandler "github.com/bing127/enterprise-agent/module-knowledge/handler/knowledge"
	knowledgeService "github.com/bing127/enterprise-agent/module-knowledge/service"
	"github.com/bing127/enterprise-agent/module-pkg/logger"
	_ "github.com/bing127/enterprise-agent/resource/docs" // swag 生成文档
	"go.uber.org/zap"
)

func main() {
	// ── 日志初始化 ────────────────────────────────────────────────
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logOutput := os.Getenv("LOG_OUTPUT") // "stdout" | "file" | "both"
	if logOutput == "" {
		logOutput = "both"
	}
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		logFilePath = "../resource/logs/app.log"
	}
	logger.InitWithConfig(logger.Config{
		Level:      logLevel,
		Env:        env,
		Output:     logOutput,
		FilePath:   logFilePath,
		MaxSizeMB:  100,
		MaxBackups: 7,
		MaxAgeDays: 30,
	})
	defer logger.Sync()

	log := logger.L().Named("main")

	// ── 配置读取 ──────────────────────────────────────────────────
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		jwtSecret = "change-me-in-production" // 仅用于开发，生产必须通过环境变量注入
	}
	jwtTTL := 2 * time.Hour

	// ── 依赖初始化 ───────────────────────────────────────────────
	// TODO: 将 nil 替换为真实存储实现（MySQL / Redis 等）
	agentSvc := agentService.NewAgentService(nil)
	knowledgeSvc := knowledgeService.NewKnowledgeService(nil)
	conversationSvc := conversationService.NewConversationService(nil, nil)

	// Auth 服务：使用内存实现（开发环境），正式接入后替换为数据库实现
	authSvc := coreService.NewAuthService(
		coreRepo.NewInMemUserRepo(),
		coreRepo.NewInMemEmailCodeRepo(),
		coreRepo.NewLogEmailSender(),
		coreService.AuthConfig{
			JWTSecret: jwtSecret,
			JWTTTL:    jwtTTL,
		},
	)

	ah := agentHandler.NewAgentHandler(agentSvc)
	kh := knowledgeHandler.NewKnowledgeHandler(knowledgeSvc)
	ch := conversationHandler.NewConversationHandler(conversationSvc)
	authH := coreHandler.NewAuthHandler(authSvc)

	// ── 启动 Hertz 服务 ──────────────────────────────────────────
	h := router.NewRouter(
		router.RouterOptions{Addr: addr, JWTSecret: jwtSecret},
		ah, kh, ch, authH,
	)

	log.Info("Enterprise Agent Gateway starting",
		zap.String("addr", addr),
		zap.String("env", env),
		zap.String("swagger", "http://localhost"+addr+"/swagger/index.html"),
	)

	h.Spin()
}
