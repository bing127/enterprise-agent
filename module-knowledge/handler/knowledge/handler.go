package knowledge

import (
	"context"

	"github.com/bing127/enterprise-agent/module-knowledge/domain/model"
	"github.com/bing127/enterprise-agent/module-knowledge/service"
	"github.com/bing127/enterprise-agent/module-pkg/logger"
	"github.com/bing127/enterprise-agent/module-pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

// KnowledgeHandler 处理 Knowledge 相关的 HTTP 请求。
type KnowledgeHandler struct {
	knowledgeService *service.KnowledgeService
	log              *zap.Logger
}

func NewKnowledgeHandler(knowledgeService *service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{
		knowledgeService: knowledgeService,
		log:              logger.L().Named("knowledge"),
	}
}

// CreateKnowledge godoc
// @Summary      创建知识库条目
// @Description  创建一个新的知识库条目
// @Tags         knowledge
// @Accept       json
// @Produce      json
// @Param        body  body      CreateKnowledgeRequest  true  "知识库信息"
// @Success      201   {object}  KnowledgeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/knowledge [post]
func (h *KnowledgeHandler) CreateKnowledge(ctx context.Context, c *app.RequestContext) {
	var req CreateKnowledgeRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.log.Warn("CreateKnowledge: bind failed", zap.Error(err))
		response.BadRequest(ctx, c, err.Error())
		return
	}
	knowledge := &model.Knowledge{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
	}
	if err := h.knowledgeService.CreateKnowledge(ctx, knowledge); err != nil {
		h.log.Error("CreateKnowledge: service error", zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	h.log.Info("CreateKnowledge: ok", zap.String("id", knowledge.ID))
	response.Created(ctx, c, toKnowledgeResponse(knowledge))
}

// GetKnowledge godoc
// @Summary      获取知识库条目
// @Description  根据 ID 获取知识库条目详情
// @Tags         knowledge
// @Produce      json
// @Param        id   path      string  true  "知识库 ID"
// @Success      200  {object}  KnowledgeResponse
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/knowledge/{id} [get]
func (h *KnowledgeHandler) GetKnowledge(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	knowledge, err := h.knowledgeService.GetKnowledgeByID(ctx, id)
	if err != nil {
		h.log.Warn("GetKnowledge: not found", zap.String("id", id), zap.Error(err))
		response.NotFound(ctx, c, err.Error())
		return
	}
	response.Success(ctx, c, toKnowledgeResponse(knowledge))
}

// ListKnowledge godoc
// @Summary      获取知识库列表
// @Description  获取所有知识库条目
// @Tags         knowledge
// @Produce      json
// @Success      200  {array}   KnowledgeResponse
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/knowledge [get]
func (h *KnowledgeHandler) ListKnowledge(ctx context.Context, c *app.RequestContext) {
	list, err := h.knowledgeService.ListKnowledge(ctx)
	if err != nil {
		h.log.Error("ListKnowledge: service error", zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	response.Success(ctx, c, toKnowledgeResponseList(list))
}

// UpdateKnowledge godoc
// @Summary      更新知识库条目
// @Description  根据 ID 更新知识库条目
// @Tags         knowledge
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "知识库 ID"
// @Param        body  body      UpdateKnowledgeRequest  true  "知识库信息"
// @Success      200   {object}  KnowledgeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/knowledge/{id} [put]
func (h *KnowledgeHandler) UpdateKnowledge(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	var req UpdateKnowledgeRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.log.Warn("UpdateKnowledge: bind failed", zap.String("id", id), zap.Error(err))
		response.BadRequest(ctx, c, err.Error())
		return
	}
	knowledge := &model.Knowledge{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
	}
	if err := h.knowledgeService.UpdateKnowledge(ctx, knowledge); err != nil {
		h.log.Error("UpdateKnowledge: service error", zap.String("id", id), zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	response.Success(ctx, c, toKnowledgeResponse(knowledge))
}

// DeleteKnowledge godoc
// @Summary      删除知识库条目
// @Description  根据 ID 删除知识库条目
// @Tags         knowledge
// @Produce      json
// @Param        id   path  string  true  "知识库 ID"
// @Success      204
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/knowledge/{id} [delete]
func (h *KnowledgeHandler) DeleteKnowledge(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if err := h.knowledgeService.DeleteKnowledge(ctx, id); err != nil {
		h.log.Error("DeleteKnowledge: service error", zap.String("id", id), zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	h.log.Info("DeleteKnowledge: ok", zap.String("id", id))
	response.Success(ctx, c, nil)
}

func toKnowledgeResponse(k *model.Knowledge) KnowledgeResponse {
	return KnowledgeResponse{
		ID:          k.ID,
		Title:       k.Title,
		Description: k.Description,
		Content:     k.Content,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
	}
}

func toKnowledgeResponseList(list []*model.Knowledge) []KnowledgeResponse {
	resp := make([]KnowledgeResponse, 0, len(list))
	for _, k := range list {
		resp = append(resp, toKnowledgeResponse(k))
	}
	return resp
}
