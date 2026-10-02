package service

import (
	"fmt"

	"github.com/idepin/pomodoro/internal/model"
	"github.com/idepin/pomodoro/internal/repository"
)

type ItemService struct {
	repo repository.ItemRepo
}

func NewItemService(repo repository.ItemRepo) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) Get(id string) (model.Item, error) {
	return s.repo.Get(id)
}

func (s *ItemService) List() []model.Item {
	return s.repo.List()
}

func (s *ItemService) Create(id, name string) (model.Item, error) {
	if id == "" || name == "" {
		return model.Item{}, fmt.Errorf("id and name required")
	}
	item := model.Item{ID: id, Name: name}
	if err := s.repo.Save(item); err != nil {
		return model.Item{}, err
	}
	return item, nil
}

func (s *ItemService) Delete(id string) error {
	return s.repo.Delete(id)
}
