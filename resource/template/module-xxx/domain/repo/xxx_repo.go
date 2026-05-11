package repo

import (
	"context"

	"github.com/bing127/enterprise-agent/module-xxx/domain/model"
)

// XxxRepo 仓储接口，由基础设施层实现。
type XxxRepo interface {
	Create(ctx context.Context, xxx *model.Xxx) error
	FindByID(ctx context.Context, id int64) (*model.Xxx, error)
	Update(ctx context.Context, xxx *model.Xxx) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, size int) ([]*model.Xxx, int64, error)
}
