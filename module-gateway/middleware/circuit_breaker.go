package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/bing127/enterprise-agent/module-pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sony/gobreaker"
)

// CircuitBreakerMiddleware 基于 sony/gobreaker 的熔断器中间件。
//
//   - name:        熔断器名称（建议用路由组或服务名）
//   - maxFailures: 连续失败次数阈值，超过后进入断开（Open）状态
//   - timeout:     断开状态持续时间，超时后进入半开（Half-Open）状态尝试恢复
//
// 熔断器将 HTTP 5xx 响应计为失败。
func CircuitBreakerMiddleware(name string, maxFailures uint32, timeout time.Duration, _ float64) app.HandlerFunc {
	st := gobreaker.Settings{
		Name:    name,
		Timeout: timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= maxFailures
		},
	}
	cb := gobreaker.NewCircuitBreaker(st)

	return func(ctx context.Context, c *app.RequestContext) {
		_, err := cb.Execute(func() (interface{}, error) {
			c.Next(ctx)
			if c.Response.StatusCode() >= 500 {
				return nil, errors.New("server error")
			}
			return nil, nil
		})

		if err == nil {
			return
		}
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			response.ServiceUnavailable(ctx, c, "service temporarily unavailable, please retry later")
			c.Abort()
		}
		// 其他错误（5xx）已由 handler 写入响应，不再重复覆盖
	}
}
