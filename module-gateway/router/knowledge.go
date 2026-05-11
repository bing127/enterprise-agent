package router

import (
	knowledgeHandler "github.com/bing127/enterprise-agent/module-knowledge/handler/knowledge"

	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterKnowledgeRoutes 注册 Knowledge 相关路由。
func RegisterKnowledgeRoutes(v1 *route.RouterGroup, kh *knowledgeHandler.KnowledgeHandler) {
	knowledge := v1.Group("/knowledge")
	{
		knowledge.POST("", kh.CreateKnowledge)
		knowledge.GET("", kh.ListKnowledge)
		knowledge.GET("/:id", kh.GetKnowledge)
		knowledge.PUT("/:id", kh.UpdateKnowledge)
		knowledge.DELETE("/:id", kh.DeleteKnowledge)
	}
}
