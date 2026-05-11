// Package constants 定义全局共享的常量，供所有子模块引用。
package constants

// ── HTTP 头字段 ───────────────────────────────────────────────────────────────

const (
	HeaderRequestID     = "X-Request-ID"
	HeaderTraceID       = "X-Trace-ID"
	HeaderAuthorization = "Authorization"
	HeaderRealIP        = "X-Real-IP"
	HeaderForwardedFor  = "X-Forwarded-For"
)

// ── Hertz RequestContext.Set / Get 键名 ────────────────────────────────────────

const (
	CtxKeyRequestID = "request_id"
	CtxKeyTraceID   = "trace_id"
	CtxKeyUserID    = "user_id"
	CtxKeyUserEmail = "user_email"
)

// ── 认证 ──────────────────────────────────────────────────────────────────────

const (
	BearerPrefix = "Bearer " // Authorization: Bearer <token>
)

// ── 邮件验证码类型 ─────────────────────────────────────────────────────────────

const (
	CodeTypeRegister   = "register"       // 注册验证码
	CodeTypeResetPwd   = "reset_password" // 重置密码验证码
	EmailCodeExpireMin = 10               // 验证码有效期（分钟）
	EmailCodeLen       = 6                // 验证码长度
)

// ── 分页 ──────────────────────────────────────────────────────────────────────

const (
	DefaultPageNum  = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// ── 用户状态 ──────────────────────────────────────────────────────────────────

const (
	UserStatusInactive = 0 // 待激活
	UserStatusActive   = 1 // 正常
	UserStatusBanned   = 2 // 封禁
)
