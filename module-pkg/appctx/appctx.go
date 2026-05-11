// Package appctx 提供全局用户上下文的存取工具。
//
// 使用示例（在 Hertz 中间件中注入，在 Handler 中读取）：
//
//	// 中间件：解析 JWT 后注入
//	ctx = appctx.WithUser(ctx, &appctx.UserInfo{UserID: 1, Email: "a@b.com", Nickname: "Alice"})
//
//	// Handler / Service：读取当前用户
//	u, ok := appctx.GetUser(ctx)
//	u  := appctx.MustGetUser(ctx) // 已认证的路由中使用
package appctx

import "context"

// UserInfo 保存当前请求已认证的用户信息。
type UserInfo struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// ctxUserKey 是包私有的 context key，防止与其他包的 key 碰撞。
type ctxUserKey struct{}

// WithUser 将 UserInfo 注入 context 并返回新 context。
func WithUser(ctx context.Context, user *UserInfo) context.Context {
	return context.WithValue(ctx, ctxUserKey{}, user)
}

// GetUser 从 context 中读取 UserInfo。
// 若未注入则返回 (nil, false)。
func GetUser(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(ctxUserKey{}).(*UserInfo)
	return u, ok && u != nil
}

// MustGetUser 从 context 中读取 UserInfo，若未注入则 panic。
// 只在确定已通过认证中间件的路由中调用。
func MustGetUser(ctx context.Context) *UserInfo {
	u, ok := GetUser(ctx)
	if !ok {
		panic("appctx: UserInfo not found in context; ensure AuthMiddleware runs first")
	}
	return u
}
