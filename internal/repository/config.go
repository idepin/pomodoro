package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/idepin/pomodoro/internal/model"
)

var ErrInvalidConfig = errors.New("invalid pomodoro config")

type ConfigRepo struct {
	path string
	mu   sync.Mutex
}

func NewConfigRepo(path string) *ConfigRepo { return &ConfigRepo{path: path} }

func (r *ConfigRepo) Get() (model.Config, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return model.DefaultConfig(), nil
	}
	if err != nil {
		return model.Config{}, err
	}
	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return model.Config{}, err
	}
	return cfg, nil
}

func (r *ConfigRepo) Save(cfg model.Config) error {
	if cfg.WorkMinutes <= 0 || cfg.ShortBreakMinutes <= 0 || cfg.LongBreakMinutes <= 0 || cfg.LongBreakInterval <= 0 {
		return ErrInvalidConfig
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return err
	}
	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path)
}
