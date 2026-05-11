package repo

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/bing127/enterprise-agent/module-core/domain/model"
	"github.com/bing127/enterprise-agent/module-core/domain/repo"
)

// InMemEmailCodeRepo is a thread-safe in-memory implementation of repo.EmailCodeRepo.
type InMemEmailCodeRepo struct {
	mu     sync.RWMutex
	codes  map[string]*model.EmailCode // key: email+":"+type
	nextID int64
}

// NewInMemEmailCodeRepo creates a new InMemEmailCodeRepo.
func NewInMemEmailCodeRepo() repo.EmailCodeRepo {
	return &InMemEmailCodeRepo{
		codes: make(map[string]*model.EmailCode),
	}
}

func key(email, codeType string) string { return email + ":" + codeType }

func (r *InMemEmailCodeRepo) Save(ctx context.Context, code *model.EmailCode) error {
	id := atomic.AddInt64(&r.nextID, 1)
	cp := *code
	cp.ID = id

	r.mu.Lock()
	r.codes[key(cp.Email, cp.Type)] = &cp
	r.mu.Unlock()
	return nil
}

func (r *InMemEmailCodeRepo) Find(ctx context.Context, email, codeType string) (*model.EmailCode, error) {
	r.mu.RLock()
	c := r.codes[key(email, codeType)]
	r.mu.RUnlock()
	if c == nil || c.Used {
		return nil, nil
	}
	cp := *c
	return &cp, nil
}

func (r *InMemEmailCodeRepo) MarkUsed(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.codes {
		if c.ID == id {
			c.Used = true
			return nil
		}
	}
	return nil
}
