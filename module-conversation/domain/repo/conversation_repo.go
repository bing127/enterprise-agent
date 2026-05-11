package repo

import (
	"context"
	"time"
)

// Conversation 是 repo 层的持久化结构（与 domain/model 解耦）。
type Conversation struct {
	ID        string
	UserID    string
	AgentID   string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Message 是 repo 层存储的单条消息。
type Message struct {
	ID        string
	Role      string
	Content   string
	Timestamp int64
}

// ConversationRepo defines the interface for conversation repository.
type ConversationRepo interface {
	Create(ctx context.Context, conversation *Conversation) error
	GetByID(ctx context.Context, id string) (*Conversation, error)
	Update(ctx context.Context, conversation *Conversation) error
	Delete(ctx context.Context, id string) error
}
