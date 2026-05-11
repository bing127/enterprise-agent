package router

import (
	conversationHandler "github.com/bing127/enterprise-agent/module-conversation/handler/conversation"

	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterConversationRoutes 注册 Conversation 相关路由。
func RegisterConversationRoutes(v1 *route.RouterGroup, ch *conversationHandler.ConversationHandler) {
	conversations := v1.Group("/conversations")
	{
		conversations.POST("", ch.CreateConversation)
		conversations.GET("/:id", ch.GetConversation)
		conversations.POST("/:id/messages", ch.SendMessage) // Eino AI 对话
	}
}
