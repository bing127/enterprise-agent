package conversation

import (
	"context"

	"github.com/bing127/enterprise-agent/module-conversation/domain/model"
	"github.com/bing127/enterprise-agent/module-conversation/service"
	"github.com/bing127/enterprise-agent/module-pkg/logger"
	"github.com/bing127/enterprise-agent/module-pkg/response"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

// ConversationHandler 处理 Conversation 相关的 HTTP 请求。
type ConversationHandler struct {
	service *service.ConversationService
	log     *zap.Logger
}

func NewConversationHandler(s *service.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		service: s,
		log:     logger.L().Named("conversation"),
	}
}

// CreateConversation godoc
// @Summary      创建会话
// @Description  创建一个新的对话会话
// @Tags         conversations
// @Accept       json
// @Produce      json
// @Param        body  body      CreateConversationRequest  true  "会话信息"
// @Success      201   {object}  ConversationResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/conversations [post]
func (h *ConversationHandler) CreateConversation(ctx context.Context, c *app.RequestContext) {
	var req CreateConversationRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.log.Warn("CreateConversation: bind failed", zap.Error(err))
		response.BadRequest(ctx, c, err.Error())
		return
	}
	conversation := &model.Conversation{
		UserID:  req.UserID,
		AgentID: req.AgentID,
	}
	if err := h.service.CreateConversation(ctx, conversation); err != nil {
		h.log.Error("CreateConversation: service error", zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	h.log.Info("CreateConversation: ok", zap.String("id", conversation.ID))
	response.Created(ctx, c, toConversationResponse(conversation))
}

// GetConversation godoc
// @Summary      获取会话
// @Description  根据 ID 获取会话详情
// @Tags         conversations
// @Produce      json
// @Param        id   path      string  true  "会话 ID"
// @Success      200  {object}  ConversationResponse
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/conversations/{id} [get]
func (h *ConversationHandler) GetConversation(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	conversation, err := h.service.GetConversation(ctx, id)
	if err != nil {
		h.log.Warn("GetConversation: not found", zap.String("id", id), zap.Error(err))
		response.NotFound(ctx, c, err.Error())
		return
	}
	response.Success(ctx, c, toConversationResponse(conversation))
}

// SendMessage godoc
// @Summary      发送消息（AI 对话）
// @Description  向指定会话发送用户消息，通过 Eino 驱动 AI 模型生成回复
// @Tags         conversations
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "会话 ID"
// @Param        body  body      SendMessageRequest  true  "消息内容"
// @Success      200   {object}  SendMessageResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/conversations/{id}/messages [post]
func (h *ConversationHandler) SendMessage(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	var req SendMessageRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.log.Warn("SendMessage: bind failed", zap.String("id", id), zap.Error(err))
		response.BadRequest(ctx, c, err.Error())
		return
	}
	reply, err := h.service.Chat(ctx, id, req.Content)
	if err != nil {
		h.log.Error("SendMessage: chat error", zap.String("id", id), zap.Error(err))
		response.InternalError(ctx, c, err.Error())
		return
	}
	h.log.Info("SendMessage: ok", zap.String("id", id))
	response.Success(ctx, c, SendMessageResponse{Reply: reply})
}

func toConversationResponse(conv *model.Conversation) ConversationResponse {
	msgs := make([]MessageResponse, 0, len(conv.Messages))
	for _, m := range conv.Messages {
		msgs = append(msgs, MessageResponse{
			ID:        m.ID,
			Role:      m.Role,
			Content:   m.Content,
			Timestamp: m.Timestamp,
		})
	}
	return ConversationResponse{
		ID:        conv.ID,
		UserID:    conv.UserID,
		AgentID:   conv.AgentID,
		Messages:  msgs,
		CreatedAt: conv.CreatedAt,
		UpdatedAt: conv.UpdatedAt,
	}
}
