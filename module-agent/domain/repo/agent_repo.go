package repo

import (
	"context"

	"github.com/bing127/enterprise-agent/module-agent/domain/model"
)

// AgentRepository defines the interface for agent repository operations.
type AgentRepository interface {
	Create(ctx context.Context, agent *model.Agent) error
	GetByID(ctx context.Context, id string) (*model.Agent, error)
	Update(ctx context.Context, agent *model.Agent) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Agent, error)
}
