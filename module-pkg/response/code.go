package response

// Code 业务错误码。
// 命名规则：HTTP 状态码 * 100 + 业务序号，0 表示成功。
type Code int

const (
	// ── 成功 ─────────────────────────────────────────
	CodeSuccess Code = 0

	// ── 客户端错误 4xx ─────────────────────────────────
	CodeBadRequest      Code = 40001 // 请求参数错误 / 校验失败
	CodeUnauthorized    Code = 40101 // 未认证 / Token 无效
	CodeTokenExpired    Code = 40102 // Token 已过期
	CodeForbidden       Code = 40301 // 无权限
	CodeNotFound        Code = 40401 // 资源不存在
	CodeConflict        Code = 40901 // 资源冲突（如邮箱已注册）
	CodeTooManyRequests Code = 42901 // 请求过于频繁

	// ── 服务端错误 5xx ─────────────────────────────────
	CodeInternalError      Code = 50001 // 服务内部错误
	CodeServiceUnavailable Code = 50301 // 服务不可用（熔断 / 降载）
)

// codeMessages 内置错误码文案。
var codeMessages = map[Code]string{
	CodeSuccess:            "ok",
	CodeBadRequest:         "bad request",
	CodeUnauthorized:       "unauthorized",
	CodeTokenExpired:       "token expired",
	CodeForbidden:          "forbidden",
	CodeNotFound:           "not found",
	CodeConflict:           "conflict",
	CodeTooManyRequests:    "too many requests",
	CodeInternalError:      "internal server error",
	CodeServiceUnavailable: "service unavailable",
}

// Message 返回错误码对应的默认文案。
func (c Code) Message() string {
	if msg, ok := codeMessages[c]; ok {
		return msg
	}
	return "unknown error"
}

// HTTPStatus 将业务错误码映射到 HTTP 状态码。
func (c Code) HTTPStatus() int {
	if c == CodeSuccess {
		return 200
	}
	// 取错误码前三位作为 HTTP 状态码
	httpCode := int(c) / 100
	if httpCode < 100 || httpCode > 599 {
		return 500
	}
	return httpCode
}
