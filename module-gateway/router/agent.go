package router

import (
	agentHandler "github.com/bing127/enterprise-agent/module-agent/handler/agent"

	"github.com/cloudwego/hertz/pkg/route"
)

// RegisterAgentRoutes 注册 Agent 相关路由。
func RegisterAgentRoutes(v1 *route.RouterGroup, ah *agentHandler.AgentHandler) {
	agents := v1.Group("/agents")
	{
		agents.POST("", ah.CreateAgent)
		agents.GET("", ah.ListAgents)
		agents.GET("/:id", ah.GetAgent)
		agents.PUT("/:id", ah.UpdateAgent)
		agents.DELETE("/:id", ah.DeleteAgent)
	}
}
