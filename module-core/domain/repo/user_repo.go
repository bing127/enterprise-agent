package repo

import (
	"context"

	"github.com/bing127/enterprise-agent/module-core/domain/model"
)

// UserRepo 用户存储接口。
type UserRepo interface {
	// Create 创建用户，返回带 ID 的用户对象。
	Create(ctx context.Context, user *model.User) (*model.User, error)

	// FindByEmail 根据邮箱查找用户，不存在时返回 (nil, nil)。
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// FindByID 根据 ID 查找用户。
	FindByID(ctx context.Context, id int64) (*model.User, error)

	// UpdateStatus 更新用户状态（激活、封禁等）。
	UpdateStatus(ctx context.Context, id int64, status int8) error
}

// EmailCodeRepo 邮箱验证码存储接口。
type EmailCodeRepo interface {
	// Save 保存验证码（覆盖同邮箱同类型的旧记录）。
	Save(ctx context.Context, code *model.EmailCode) error

	// Find 查找未过期且未使用的验证码。
	Find(ctx context.Context, email, codeType string) (*model.EmailCode, error)

	// MarkUsed 标记验证码为已使用。
	MarkUsed(ctx context.Context, id int64) error
}

// EmailSender 邮件发送接口，底层可以是 SMTP / SES / 第三方服务。
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}
