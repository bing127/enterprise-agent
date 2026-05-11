package model

import "time"

// User 系统用户。
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // bcrypt 散列，不序列化到 JSON
	Nickname     string    `json:"nickname"`
	Status       int8      `json:"status"` // 0=待激活 1=正常 2=封禁
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EmailCode 邮箱验证码，用于注册或重置密码。
type EmailCode struct {
	ID        int64
	Email     string
	Code      string
	Type      string // "register" | "reset_password"
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}
