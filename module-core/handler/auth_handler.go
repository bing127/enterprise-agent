package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/bing127/enterprise-agent/module-core/service"
	"github.com/bing127/enterprise-agent/module-pkg/appctx"
	"github.com/bing127/enterprise-agent/module-pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// AuthHandler HTTP 认证处理器。
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler 创建 AuthHandler 实例。
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// ── 请求体 DTO ────────────────────────────────────────────────────────────────

type sendCodeBody struct {
	Email    string `json:"email"`
	CodeType string `json:"code_type"`
}

type registerBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Nickname string `json:"nickname"`
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ── 处理器方法 ─────────────────────────────────────────────────────────────────

// SendCode 发送邮箱验证码。
//
//	@Summary     发送邮箱验证码
//	@Description 向指定邮箱发送 6 位数字验证码，有效期 10 分钟
//	@Tags        auth
//	@Accept      json
//	@Produce     json
//	@Param       body body sendCodeBody true "请求体"
//	@Success     200  {object} response.R
//	@Failure     400  {object} response.R
//	@Router      /api/v1/auth/send-code [post]
func (h *AuthHandler) SendCode(ctx context.Context, c *app.RequestContext) {
	var req sendCodeBody
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(ctx, c, "invalid request body")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		response.BadRequest(ctx, c, "email is required")
		return
	}
	if req.CodeType == "" {
		req.CodeType = "register"
	}

	if err := h.svc.SendEmailCode(ctx, &service.SendCodeRequest{
		Email:    req.Email,
		CodeType: req.CodeType,
	}); err != nil {
		response.InternalError(ctx, c, err.Error())
		return
	}
	response.Success(ctx, c, nil)
}

// Register 邮箱注册。
//
//	@Summary     邮箱注册
//	@Description 验证邮箱验证码后创建用户账号
//	@Tags        auth
//	@Accept      json
//	@Produce     json
//	@Param       body body registerBody true "请求体"
//	@Success     200  {object} response.R
//	@Failure     400  {object} response.R
//	@Failure     409  {object} response.R
//	@Router      /api/v1/auth/register [post]
func (h *AuthHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req registerBody
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(ctx, c, "invalid request body")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" || req.Code == "" {
		response.BadRequest(ctx, c, "email, password and code are required")
		return
	}

	err := h.svc.Register(ctx, &service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Code:     req.Code,
		Nickname: req.Nickname,
	})
	if err != nil {
		switch err {
		case service.ErrEmailExists:
			response.Conflict(ctx, c, "email already registered")
		case service.ErrInvalidCode:
			response.BadRequest(ctx, c, "invalid or expired verification code")
		default:
			response.InternalError(ctx, c, err.Error())
		}
		return
	}
	response.Success(ctx, c, nil)
}

// Login 邮箱密码登录。
//
//	@Summary     用户登录
//	@Description 邮箱 + 密码登录，返回 JWT 令牌
//	@Tags        auth
//	@Accept      json
//	@Produce     json
//	@Param       body body loginBody true "请求体"
//	@Success     200  {object} response.R{data=service.LoginResponse}
//	@Failure     400  {object} response.R
//	@Failure     401  {object} response.R
//	@Router      /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req loginBody
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(ctx, c, "invalid request body")
		return
	}
	if req.Email == "" || req.Password == "" {
		response.BadRequest(ctx, c, "email and password are required")
		return
	}

	resp, err := h.svc.Login(ctx, &service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrWrongPassword):
			response.Unauthorized(ctx, c, "invalid email or password")
		case errors.Is(err, service.ErrUserNotActive):
			response.Forbidden(ctx, c, "account is not active")
		default:
			response.InternalError(ctx, c, err.Error())
		}
		return
	}
	response.Success(ctx, c, resp)
}

// Me 获取当前登录用户信息。
//
//	@Summary     获取当前用户信息
//	@Description 从 JWT 中解析并返回当前登录用户的基本信息
//	@Tags        auth
//	@Produce     json
//	@Security    BearerAuth
//	@Success     200 {object} response.R{data=appctx.UserInfo}
//	@Failure     401 {object} response.R
//	@Router      /api/v1/auth/me [get]
func (h *AuthHandler) Me(ctx context.Context, c *app.RequestContext) {
	u, ok := appctx.GetUser(ctx)
	if !ok {
		response.Unauthorized(ctx, c, "not authenticated")
		return
	}
	response.Success(ctx, c, u)
}
