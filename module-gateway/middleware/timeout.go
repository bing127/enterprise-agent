package middleware

import (
	"context"
	"time"

	"github.com/bing127/enterprise-agent/module-pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// TimeoutMiddleware 为每个请求设置最大处理时间。
//
// 若 handler 在 d 时间内未返回，返回 503 并中断后续处理。
// 注意：Hertz handler 中若有阻塞 I/O，应使用 ctx.Done() 监听取消信号以提前退出。
func TimeoutMiddleware(d time.Duration) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		timeoutCtx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		done := make(chan struct{}, 1)
		go func() {
			c.Next(timeoutCtx)
			done <- struct{}{}
		}()

		select {
		case <-done:
			// 正常完成
		case <-timeoutCtx.Done():
			response.ServiceUnavailable(timeoutCtx, c, "request timeout")
			c.Abort()
		}
	}
}
