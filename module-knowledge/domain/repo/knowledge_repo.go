package repo

import (
	"context"

	"github.com/bing127/enterprise-agent/module-knowledge/domain/model"
)

// KnowledgeRepository defines the interface for knowledge repository operations.
type KnowledgeRepository interface {
	Create(ctx context.Context, knowledge *model.Knowledge) error
	GetByID(ctx context.Context, id string) (*model.Knowledge, error)
	Update(ctx context.Context, knowledge *model.Knowledge) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Knowledge, error)
}
