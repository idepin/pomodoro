package repository

import (
	"errors"
	"sync"

	"github.com/idepin/pomodoro/internal/model"
)

var ErrNotFound = errors.New("item not found")

// ItemRepo abstracts storage so service never touches infra directly.
type ItemRepo interface {
	Get(id string) (model.Item, error)
	List() []model.Item
	Save(item model.Item) error
	Delete(id string) error
}

type MemoryItemRepo struct {
	mu    sync.RWMutex
	items map[string]model.Item
}

func NewMemoryItemRepo() *MemoryItemRepo {
	return &MemoryItemRepo{items: make(map[string]model.Item)}
}

func (r *MemoryItemRepo) Get(id string) (model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok {
		return model.Item{}, ErrNotFound
	}
	return item, nil
}

func (r *MemoryItemRepo) List() []model.Item {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Item, 0, len(r.items))
	for _, it := range r.items {
		out = append(out, it)
	}
	return out
}

func (r *MemoryItemRepo) Save(item model.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[item.ID] = item
	return nil
}

func (r *MemoryItemRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
