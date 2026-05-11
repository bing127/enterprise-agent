package conversation

// CreateConversationRequest 创建会话的请求体。
type CreateConversationRequest struct {
	UserID  string `json:"user_id"  vd:"len($)>0"`
	AgentID string `json:"agent_id" vd:"len($)>0"`
}

// SendMessageRequest 用户发送消息的请求体。
type SendMessageRequest struct {
	Content string `json:"content" vd:"len($)>0"`
}
