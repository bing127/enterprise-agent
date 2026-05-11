// Package response 提供统一的 HTTP 响应格式。
//
// 所有接口返回结构：
//
//	{
//	  "code":    0,     // 业务错误码，0 表示成功
//	  "success": true,  // 是否成功
//	  "msg":     "ok",  // 描述信息
//	  "data":    {...}  // 业务数据，失败时为 null
//	}
package response

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

// R 标准响应体。
type R struct {
	Code    Code        `json:"code"`
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// Success 返回成功响应（HTTP 200）。
func Success(ctx context.Context, c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusOK, R{
		Code:    CodeSuccess,
		Success: true,
		Msg:     CodeSuccess.Message(),
		Data:    data,
	})
}

// Created 返回创建成功响应（HTTP 201）。
func Created(ctx context.Context, c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusCreated, R{
		Code:    CodeSuccess,
		Success: true,
		Msg:     CodeSuccess.Message(),
		Data:    data,
	})
}

// Fail 返回失败响应，HTTP 状态码由 code 自动推断。
func Fail(ctx context.Context, c *app.RequestContext, code Code, msg string) {
	if msg == "" {
		msg = code.Message()
	}
	c.JSON(code.HTTPStatus(), R{
		Code:    code,
		Success: false,
		Msg:     msg,
		Data:    nil,
	})
}

// BadRequest 400 参数校验失败。
func BadRequest(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeBadRequest, msg)
}

// Unauthorized 401 未认证。
func Unauthorized(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeUnauthorized, msg)
}

// Forbidden 403 无权限。
func Forbidden(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeForbidden, msg)
}

// NotFound 404 资源不存在。
func NotFound(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeNotFound, msg)
}

// Conflict 409 资源冲突。
func Conflict(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeConflict, msg)
}

// TooManyRequests 429 请求过于频繁。
func TooManyRequests(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeTooManyRequests, msg)
}

// InternalError 500 服务内部错误。
func InternalError(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeInternalError, msg)
}

// ServiceUnavailable 503 服务不可用（熔断 / 降载）。
func ServiceUnavailable(ctx context.Context, c *app.RequestContext, msg string) {
	Fail(ctx, c, CodeServiceUnavailable, msg)
}
