package agent

// CreateAgentRequest 创建 Agent 的请求体。
type CreateAgentRequest struct {
	Name        string `json:"name"        vd:"len($)>0"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// UpdateAgentRequest 更新 Agent 的请求体。
type UpdateAgentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
