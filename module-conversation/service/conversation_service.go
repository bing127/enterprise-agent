package service

import (
	"context"

	"github.com/bing127/enterprise-agent/module-conversation/domain/model"
	"github.com/bing127/enterprise-agent/module-conversation/domain/repo"

	einoModel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// ConversationService 管理会话，并通过 Eino ChatModel 驱动 AI 对话。
type ConversationService struct {
	conversationRepo repo.ConversationRepo
	chatModel        einoModel.ChatModel // Eino 抽象 ChatModel，可对接 OpenAI / Ark 等
}

// NewConversationService 创建 ConversationService。
// chatModel 可为 nil（开发/测试时跳过 AI 调用）。
func NewConversationService(conversationRepo repo.ConversationRepo, chatModel einoModel.ChatModel) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		chatModel:        chatModel,
	}
}

func (s *ConversationService) CreateConversation(ctx context.Context, conversation *model.Conversation) error {
	c := &repo.Conversation{
		ID:      conversation.ID,
		UserID:  conversation.UserID,
		AgentID: conversation.AgentID,
	}
	return s.conversationRepo.Create(ctx, c)
}

func (s *ConversationService) GetConversation(ctx context.Context, id string) (*model.Conversation, error) {
	c, err := s.conversationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	msgs := make([]model.Message, 0, len(c.Messages))
	for _, m := range c.Messages {
		msgs = append(msgs, model.Message{
			ID:        m.ID,
			Role:      m.Role,
			Content:   m.Content,
			Timestamp: m.Timestamp,
		})
	}
	return &model.Conversation{
		ID:        c.ID,
		UserID:    c.UserID,
		AgentID:   c.AgentID,
		Messages:  msgs,
		CreatedAt: c.CreatedAt.UnixMilli(),
		UpdatedAt: c.UpdatedAt.UnixMilli(),
	}, nil
}

// Chat 将用户消息通过 Eino 发送给 AI 模型，返回 AI 回复内容。
// 若未配置 chatModel，返回占位回复，便于本地开发调试。
func (s *ConversationService) Chat(ctx context.Context, conversationID, userInput string) (string, error) {
	if s.chatModel == nil {
		return "[AI 未配置，请在 main.go 中注入 Eino ChatModel]", nil
	}

	// 构建对话消息列表：系统提示 + 用户输入
	messages := []*schema.Message{
		schema.SystemMessage("你是一位专业的企业级智能助手，请简洁、准确地回答用户问题。"),
		schema.UserMessage(userInput),
	}

	resp, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}
