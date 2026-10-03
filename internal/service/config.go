package service

import (
	"github.com/idepin/pomodoro/internal/model"
	"github.com/idepin/pomodoro/internal/repository"
)

type ConfigService struct {
	repo *repository.ConfigRepo
}

func NewConfigService(repo *repository.ConfigRepo) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) Get() (model.Config, error)  { return s.repo.Get() }
func (s *ConfigService) Save(cfg model.Config) error { return s.repo.Save(cfg) }
