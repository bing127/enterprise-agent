package model

type Conversation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	AgentID   string    `json:"agent_id"`
	Messages  []Message `json:"messages"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type Message struct {
	ID        string `json:"id"`
	Role      string `json:"role"` // user / assistant / system
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
