package conversation

// MessageResponse 单条消息响应体。
type MessageResponse struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// ConversationResponse 会话响应体。
type ConversationResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	AgentID   string            `json:"agent_id"`
	Messages  []MessageResponse `json:"messages"`
	CreatedAt int64             `json:"created_at"`
	UpdatedAt int64             `json:"updated_at"`
}

// SendMessageResponse AI 回复的响应体。
type SendMessageResponse struct {
	Reply string `json:"reply"`
}
