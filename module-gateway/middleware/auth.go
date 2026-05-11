package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/bing127/enterprise-agent/module-pkg/appctx"
	pkgjwt "github.com/bing127/enterprise-agent/module-pkg/jwt"
	"github.com/bing127/enterprise-agent/module-pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// AuthMiddleware 验证 Bearer JWT 令牌，并将用户信息注入 context。
//
//   - jwtSecret: JWT 签名密钥（从配置读取后传入）
func AuthMiddleware(jwtSecret string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			response.Unauthorized(ctx, c, "missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") || parts[1] == "" {
			response.Unauthorized(ctx, c, "invalid Authorization format, expected: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := pkgjwt.Parse(jwtSecret, parts[1])
		if err != nil {
			if errors.Is(err, pkgjwt.ErrExpiredToken) {
				response.Fail(ctx, c, response.CodeTokenExpired, "token expired")
			} else {
				response.Unauthorized(ctx, c, "invalid token")
			}
			c.Abort()
			return
		}

		// 将用户信息注入 context，后续 handler/service 可通过 appctx.GetUser(ctx) 获取
		ctx = appctx.WithUser(ctx, &appctx.UserInfo{
			UserID:   claims.UserID,
			Email:    claims.Email,
			Nickname: claims.Nickname,
		})

		c.Next(ctx)
	}
}
