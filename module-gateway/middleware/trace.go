package middleware

import (
	"context"
	"time"

	"github.com/bing127/enterprise-agent/module-pkg/constants"
	"github.com/bing127/enterprise-agent/module-pkg/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TraceMiddleware 注入 X-Request-ID 并记录每次请求的结构化日志。
//
// X-Request-ID 优先使用客户端传入的值，否则自动生成 UUID v4。
// 请求 ID 同时存入 Hertz RequestContext（key: constants.CtxKeyRequestID）。
func TraceMiddleware() app.HandlerFunc {
	log := logger.L().Named("access")
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 确定 Request ID
		reqID := string(c.GetHeader(constants.HeaderRequestID))
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set(constants.CtxKeyRequestID, reqID)
		c.Header(constants.HeaderRequestID, reqID) // 回写到响应头

		start := time.Now()
		c.Next(ctx)

		log.Info("request",
			zap.String("request_id", reqID),
			zap.ByteString("method", c.Method()),
			zap.ByteString("path", c.Path()),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status", c.Response.StatusCode()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

// CORSMiddleware 允许跨域请求。
func CORSMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next(ctx)
	}
}
