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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	agentHandler "github.com/bing127/enterprise-agent/module-agent/handler/agent"
	agentService "github.com/bing127/enterprise-agent/module-agent/service"
	moduleconfig "github.com/bing127/enterprise-agent/module-config"
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
	cfg, err := loadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load config failed: %v\n", err)
		os.Exit(1)
	}

	// ── 日志初始化 ────────────────────────────────────────────────
	env := cfg.App.Env
	if env == "" {
		env = "dev"
	}
	logLevel := cfg.Log.Level
	if logLevel == "" {
		logLevel = "info"
	}
	logOutput := cfg.Log.Output
	if logOutput == "" {
		logOutput = "stdout"
	}
	logFilePath := cfg.Log.FilePath
	if logFilePath == "" {
		logFilePath = "./resource/logs/app.log"
	}
	maxSizeMB := cfg.Log.MaxSize
	if maxSizeMB <= 0 {
		maxSizeMB = 100
	}
	maxBackups := cfg.Log.MaxBackups
	if maxBackups <= 0 {
		maxBackups = 7
	}
	maxAgeDays := cfg.Log.MaxAge
	if maxAgeDays <= 0 {
		maxAgeDays = 30
	}

	logger.InitWithConfig(logger.Config{
		Level:      logLevel,
		Env:        env,
		Output:     logOutput,
		FilePath:   logFilePath,
		MaxSizeMB:  maxSizeMB,
		MaxBackups: maxBackups,
		MaxAgeDays: maxAgeDays,
	})
	defer logger.Sync()

	log := logger.L().Named("main")

	// ── 配置读取 ──────────────────────────────────────────────────
	host := cfg.Server.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}
	addr := fmt.Sprintf("%s:%d", host, port)

	jwtSecret := cfg.Auth.SecretKey
	if jwtSecret == "" {
		jwtSecret = "change-me-in-production"
	}
	jwtTTL := time.Duration(cfg.Auth.TokenTTL) * time.Second
	if jwtTTL <= 0 {
		jwtTTL = 2 * time.Hour
	}

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
		zap.String("swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", port)),
	)

	h.Spin()
}

func loadConfig() (*moduleconfig.Config, error) {
	configDir := strings.TrimSpace(os.Getenv("CONFIG_DIR"))
	if configDir == "" {
		for _, candidate := range []string{"manifest/config", "../manifest/config"} {
			if st, err := os.Stat(candidate); err == nil && st.IsDir() {
				configDir = candidate
				break
			}
		}
	}
	if configDir == "" {
		return nil, fmt.Errorf("config directory not found, set CONFIG_DIR explicitly")
	}

	configDir, err := filepath.Abs(configDir)
	if err != nil {
		return nil, fmt.Errorf("resolve config dir: %w", err)
	}
	return moduleconfig.Load(configDir)
}
