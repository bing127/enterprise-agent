package service

import (
	"context"
	"errors"

	"github.com/bing127/enterprise-agent/module-knowledge/domain/model"
	"github.com/bing127/enterprise-agent/module-knowledge/domain/repo"
)

type KnowledgeService struct {
	knowledgeRepo repo.KnowledgeRepository
}

func NewKnowledgeService(knowledgeRepo repo.KnowledgeRepository) *KnowledgeService {
	return &KnowledgeService{
		knowledgeRepo: knowledgeRepo,
	}
}

func (ks *KnowledgeService) GetKnowledgeByID(ctx context.Context, id string) (*model.Knowledge, error) {
	knowledge, err := ks.knowledgeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if knowledge == nil {
		return nil, errors.New("knowledge not found")
	}
	return knowledge, nil
}

func (ks *KnowledgeService) ListKnowledge(ctx context.Context) ([]*model.Knowledge, error) {
	return ks.knowledgeRepo.List(ctx)
}

func (ks *KnowledgeService) CreateKnowledge(ctx context.Context, knowledge *model.Knowledge) error {
	return ks.knowledgeRepo.Create(ctx, knowledge)
}

func (ks *KnowledgeService) UpdateKnowledge(ctx context.Context, knowledge *model.Knowledge) error {
	return ks.knowledgeRepo.Update(ctx, knowledge)
}

func (ks *KnowledgeService) DeleteKnowledge(ctx context.Context, id string) error {
	return ks.knowledgeRepo.Delete(ctx, id)
}
