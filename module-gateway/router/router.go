package router

import (
	"context"
	"net/http"
	"time"

	agentHandler "github.com/bing127/enterprise-agent/module-agent/handler/agent"
	conversationHandler "github.com/bing127/enterprise-agent/module-conversation/handler/conversation"
	coreHandler "github.com/bing127/enterprise-agent/module-core/handler"
	"github.com/bing127/enterprise-agent/module-gateway/middleware"
	knowledgeHandler "github.com/bing127/enterprise-agent/module-knowledge/handler/knowledge"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzsSwagger "github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// RouterOptions 路由启动配置。
type RouterOptions struct {
	Addr      string // 监听地址，如 ":8080"
	JWTSecret string // JWT 签名密钥
}

// NewRouter 初始化 Hertz 服务器，注册全局中间件和所有业务路由。
func NewRouter(
	opts RouterOptions,
	ah *agentHandler.AgentHandler,
	kh *knowledgeHandler.KnowledgeHandler,
	ch *conversationHandler.ConversationHandler,
	authH *coreHandler.AuthHandler,
) *server.Hertz {
	h := server.Default(server.WithHostPorts(opts.Addr))

	// ── 全局中间件（按顺序执行）────────────────────────────────
	h.Use(
		middleware.TimeoutMiddleware(30*time.Second),   // 超时控制
		middleware.TraceMiddleware(),                    // 链路追踪 / Request ID
		middleware.CORSMiddleware(),                     // 跨域
		middleware.AdaptiveSheddingMiddleware(middleware.SheddingOptions{ // 自适应降载
			MaxCPU:      75,
			MaxInFlight: 500,
		}),
		middleware.RateLimitMiddleware(200, time.Minute), // IP 级全局限流
	)

	// ── 基础路由 ──────────────────────────────────────────────
	h.GET("/swagger/*any", hertzsSwagger.WrapHandler(swaggerFiles.Handler))
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// ── API v1 ────────────────────────────────────────────────
	v1 := h.Group("/api/v1")

	// 认证路由（无需 JWT）
	authGroup := v1.Group("/auth")
	authGroup.Use(middleware.CircuitBreakerMiddleware("auth", 5, 30*time.Second, 0.6))
	{
		authGroup.POST("/send-code", authH.SendCode)
		authGroup.POST("/register", authH.Register)
		authGroup.POST("/login", authH.Login)
	}

	// 需要 JWT 认证的路由
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(opts.JWTSecret))
	protected.Use(middleware.CircuitBreakerMiddleware("api", 10, 60*time.Second, 0.6))
	{
		protected.GET("/auth/me", authH.Me)
		RegisterAgentRoutes(protected, ah)
		RegisterKnowledgeRoutes(protected, kh)
		RegisterConversationRoutes(protected, ch)
	}

	return h
}
