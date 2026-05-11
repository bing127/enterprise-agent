package service

import (
	"context"

	"github.com/bing127/enterprise-agent/module-agent/domain/model"
	"github.com/bing127/enterprise-agent/module-agent/domain/repo"
)

type AgentService struct {
	agentRepo repo.AgentRepository
}

func NewAgentService(agentRepo repo.AgentRepository) *AgentService {
	return &AgentService{
		agentRepo: agentRepo,
	}
}

func (s *AgentService) CreateAgent(ctx context.Context, agent *model.Agent) error {
	return s.agentRepo.Create(ctx, agent)
}

func (s *AgentService) GetAgent(ctx context.Context, id string) (*model.Agent, error) {
	return s.agentRepo.GetByID(ctx, id)
}

func (s *AgentService) UpdateAgent(ctx context.Context, agent *model.Agent) error {
	return s.agentRepo.Update(ctx, agent)
}

func (s *AgentService) DeleteAgent(ctx context.Context, id string) error {
	return s.agentRepo.Delete(ctx, id)
}

func (s *AgentService) ListAgents(ctx context.Context) ([]*model.Agent, error) {
	return s.agentRepo.List(ctx)
}
