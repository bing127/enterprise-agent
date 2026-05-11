package knowledge

// CreateKnowledgeRequest 创建知识库条目的请求体。
type CreateKnowledgeRequest struct {
	Title       string `json:"title"       vd:"len($)>0"`
	Description string `json:"description"`
	Content     string `json:"content"     vd:"len($)>0"`
}

// UpdateKnowledgeRequest 更新知识库条目的请求体。
type UpdateKnowledgeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}
