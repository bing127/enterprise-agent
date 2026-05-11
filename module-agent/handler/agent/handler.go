package agent

import (
	"context"

	"github.com/bing127/enterprise-agent/module-agent/domain/model"
	"github.com/bing127/enterprise-agent/module-agent/service"
	"github.com/bing127/enterprise-agent/module-pkg/logger"
	"github.com/bing127/enterprise-agent/module-pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

// AgentHandler 处理 Agent 相关的 HTTP 请求。
type AgentHandler struct {
	agentService *service.AgentService
	log          *zap.Logger
}

func NewAgentHandler(agentService *service.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
		log:          logger.L().Named("agent"),
	}
}

// CreateAgent godoc
// @Summary      创建 Agent
// @Description  创建一个新的企业智能 Agent
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        body  body      CreateAgentRequest  true  "Agent 信息"
// @Success      201   {object}  AgentResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/agents [post]
func (h *AgentHandler) CreateAgent(ctx context.Context, c *app.RequestContext) {
	var req CreateAgentRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.log.Warn("CreateAgent: bind failed", zap.Error(err))
		response.BadRequest(ctx, c, err.Error())
		return
	}
	agent := &model.Agent{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	if err := h.agentService.CreateAgent(ctx, agent); err != nil {
		h.log.Error("CreateAgent: service error", zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	h.log.Info("CreateAgent: ok", zap.String("id", agent.ID))
	response.Created(ctx, c, toAgentResponse(agent))
}

// GetAgent godoc
// @Summary      获取 Agent
// @Description  根据 ID 获取 Agent 详情
// @Tags         agents
// @Produce      json
// @Param        id   path      string  true  "Agent ID"
// @Success      200  {object}  AgentResponse
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/agents/{id} [get]
func (h *AgentHandler) GetAgent(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	agent, err := h.agentService.GetAgent(ctx, id)
	if err != nil {
		h.log.Warn("GetAgent: not found", zap.String("id", id), zap.Error(err))
		response.NotFound(ctx, c, "agent not found")
		return
	}
	response.Success(ctx, c, toAgentResponse(agent))
}

// ListAgents godoc
// @Summary      获取 Agent 列表
// @Description  获取所有 Agent
// @Tags         agents
// @Produce      json
// @Success      200  {array}   AgentResponse
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/agents [get]
func (h *AgentHandler) ListAgents(ctx context.Context, c *app.RequestContext) {
	agents, err := h.agentService.ListAgents(ctx)
	if err != nil {
		h.log.Error("ListAgents: service error", zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	response.Success(ctx, c, toAgentResponseList(agents))
}

// UpdateAgent godoc
// @Summary      更新 Agent
// @Description  根据 ID 更新 Agent 信息
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "Agent ID"
// @Param        body  body      UpdateAgentRequest  true  "Agent 信息"
// @Success      200   {object}  AgentResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/agents/{id} [put]
func (h *AgentHandler) UpdateAgent(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	var req UpdateAgentRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.log.Warn("UpdateAgent: bind failed", zap.String("id", id), zap.Error(err))
		response.BadRequest(ctx, c, err.Error())
		return
	}
	agent := &model.Agent{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	if err := h.agentService.UpdateAgent(ctx, agent); err != nil {
		h.log.Error("UpdateAgent: service error", zap.String("id", id), zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	response.Success(ctx, c, toAgentResponse(agent))
}

// DeleteAgent godoc
// @Summary      删除 Agent
// @Description  根据 ID 删除 Agent
// @Tags         agents
// @Produce      json
// @Param        id   path  string  true  "Agent ID"
// @Success      204
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/agents/{id} [delete]
func (h *AgentHandler) DeleteAgent(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if err := h.agentService.DeleteAgent(ctx, id); err != nil {
		h.log.Error("DeleteAgent: service error", zap.String("id", id), zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	h.log.Info("DeleteAgent: ok", zap.String("id", id))
	response.Success(ctx, c, nil)
}

func toAgentResponse(a *model.Agent) AgentResponse {
	return AgentResponse{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Status:      a.Status,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func toAgentResponseList(list []*model.Agent) []AgentResponse {
	resp := make([]AgentResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, toAgentResponse(a))
	}
	return resp
}
