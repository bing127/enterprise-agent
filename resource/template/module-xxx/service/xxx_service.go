package service

import (
	"context"

	"github.com/bing127/enterprise-agent/module-xxx/domain/model"
	"github.com/bing127/enterprise-agent/module-xxx/domain/repo"
)

// XxxService 业务接口。
type XxxService interface {
	Create(ctx context.Context, xxx *model.Xxx) error
	GetByID(ctx context.Context, id int64) (*model.Xxx, error)
	Update(ctx context.Context, xxx *model.Xxx) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, size int) ([]*model.Xxx, int64, error)
}

type xxxService struct {
	repo repo.XxxRepo
}

// NewXxxService 创建 XxxService 实例。
func NewXxxService(repo repo.XxxRepo) XxxService {
	return &xxxService{repo: repo}
}

func (s *xxxService) Create(ctx context.Context, xxx *model.Xxx) error {
	return s.repo.Create(ctx, xxx)
}

func (s *xxxService) GetByID(ctx context.Context, id int64) (*model.Xxx, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *xxxService) Update(ctx context.Context, xxx *model.Xxx) error {
	return s.repo.Update(ctx, xxx)
}

func (s *xxxService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *xxxService) List(ctx context.Context, page, size int) ([]*model.Xxx, int64, error) {
	return s.repo.List(ctx, page, size)
}
