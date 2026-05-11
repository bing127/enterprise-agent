package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/bing127/enterprise-agent/module-core/domain/model"
	"github.com/bing127/enterprise-agent/module-core/domain/repo"
	"github.com/bing127/enterprise-agent/module-pkg/constants"
	pkgjwt "github.com/bing127/enterprise-agent/module-pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// ── 请求 / 响应 DTO ───────────────────────────────────────────────────────────

// SendCodeRequest 发送验证码请求。
type SendCodeRequest struct {
	Email    string `json:"email"`
	CodeType string `json:"code_type"` // "register" | "reset_password"
}

// RegisterRequest 邮箱注册请求。
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse 登录成功响应。
type LoginResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
	User     *UserDTO  `json:"user"`
}

// UserDTO 用于对外返回的用户信息（隐藏敏感字段）。
type UserDTO struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Status   int8   `json:"status"`
}

// ── 错误定义 ──────────────────────────────────────────────────────────────────

var (
	ErrEmailExists   = errors.New("email already registered")
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrInvalidCode   = errors.New("invalid or expired verification code")
	ErrUserNotActive = errors.New("user account is not active")
)

// ── 接口定义 ──────────────────────────────────────────────────────────────────

// AuthService 用户认证业务接口。
type AuthService interface {
	// SendEmailCode 发送邮箱验证码。
	SendEmailCode(ctx context.Context, req *SendCodeRequest) error
	// Register 邮箱注册（需先验证邮箱验证码）。
	Register(ctx context.Context, req *RegisterRequest) error
	// Login 邮箱密码登录，返回 JWT 令牌。
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

// ── 实现 ──────────────────────────────────────────────────────────────────────

type authService struct {
	userRepo  repo.UserRepo
	codeRepo  repo.EmailCodeRepo
	emailer   repo.EmailSender
	jwtSecret string
	jwtTTL    time.Duration
}

// AuthConfig 认证服务配置，供 Wire 注入使用。
type AuthConfig struct {
	JWTSecret string
	JWTTTL    time.Duration
}

// NewAuthService 创建 AuthService 实例。
func NewAuthService(
	userRepo repo.UserRepo,
	codeRepo repo.EmailCodeRepo,
	emailer repo.EmailSender,
	cfg AuthConfig,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		codeRepo:  codeRepo,
		emailer:   emailer,
		jwtSecret: cfg.JWTSecret,
		jwtTTL:    cfg.JWTTTL,
	}
}

// SendEmailCode 生成并发送邮箱验证码。
func (s *authService) SendEmailCode(ctx context.Context, req *SendCodeRequest) error {
	code, err := generateCode(constants.EmailCodeLen)
	if err != nil {
		return err
	}

	ec := &model.EmailCode{
		Email:     req.Email,
		Code:      code,
		Type:      req.CodeType,
		ExpiresAt: time.Now().Add(constants.EmailCodeExpireMin * time.Minute),
		Used:      false,
		CreatedAt: time.Now(),
	}
	if err := s.codeRepo.Save(ctx, ec); err != nil {
		return err
	}

	subject := "Enterprise Agent 验证码"
	body := fmt.Sprintf("您的验证码为：%s，%d 分钟内有效。", code, constants.EmailCodeExpireMin)
	return s.emailer.Send(ctx, req.Email, subject, body)
}

// Register 验证验证码后创建用户。
func (s *authService) Register(ctx context.Context, req *RegisterRequest) error {
	// 1. 检查邮箱是否已注册
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrEmailExists
	}

	// 2. 校验验证码
	if err := s.verifyCode(ctx, req.Email, req.Code, constants.CodeTypeRegister); err != nil {
		return err
	}

	// 3. 密码散列
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. 创建用户
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Email
	}
	_, err = s.userRepo.Create(ctx, &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Nickname:     nickname,
		Status:       constants.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	return err
}

// Login 邮箱密码登录。
func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.Status != constants.UserStatusActive {
		return nil, ErrUserNotActive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrWrongPassword
	}

	token, err := pkgjwt.Sign(s.jwtSecret, s.jwtTTL, user.ID, user.Email, user.Nickname)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:    token,
		ExpireAt: time.Now().Add(s.jwtTTL),
		User: &UserDTO{
			ID:       user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
			Status:   user.Status,
		},
	}, nil
}

// ── 私有辅助 ──────────────────────────────────────────────────────────────────

func (s *authService) verifyCode(ctx context.Context, email, code, codeType string) error {
	ec, err := s.codeRepo.Find(ctx, email, codeType)
	if err != nil {
		return err
	}
	if ec == nil || ec.Code != code || time.Now().After(ec.ExpiresAt) {
		return ErrInvalidCode
	}
	return s.codeRepo.MarkUsed(ctx, ec.ID)
}

// generateCode 生成指定长度的纯数字随机验证码（使用 crypto/rand）。
func generateCode(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[n.Int64()]
	}
	return string(result), nil
}
