// Package repo provides in-memory implementations of the repository interfaces.
// These are intended for development and testing only.
package repo

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/bing127/enterprise-agent/module-core/domain/model"
	"github.com/bing127/enterprise-agent/module-core/domain/repo"
)

// InMemUserRepo is a thread-safe in-memory implementation of repo.UserRepo.
type InMemUserRepo struct {
	mu      sync.RWMutex
	byEmail map[string]*model.User
	byID    map[int64]*model.User
	nextID  int64
}

// NewInMemUserRepo creates a new InMemUserRepo.
func NewInMemUserRepo() repo.UserRepo {
	return &InMemUserRepo{
		byEmail: make(map[string]*model.User),
		byID:    make(map[int64]*model.User),
		nextID:  0,
	}
}

func (r *InMemUserRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	id := atomic.AddInt64(&r.nextID, 1)
	cp := *user
	cp.ID = id

	r.mu.Lock()
	r.byEmail[cp.Email] = &cp
	r.byID[cp.ID] = &cp
	r.mu.Unlock()

	return &cp, nil
}

func (r *InMemUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	u := r.byEmail[email]
	r.mu.RUnlock()
	if u == nil {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}

func (r *InMemUserRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	r.mu.RLock()
	u := r.byID[id]
	r.mu.RUnlock()
	if u == nil {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}

func (r *InMemUserRepo) UpdateStatus(ctx context.Context, id int64, status int8) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return nil
	}
	u.Status = status
	r.byEmail[u.Email].Status = status
	return nil
}
