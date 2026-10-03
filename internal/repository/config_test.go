package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/idepin/pomodoro/internal/model"
)

func TestConfigRepo_GetDefault(t *testing.T) {
	dir := t.TempDir()
	repo := NewConfigRepo(filepath.Join(dir, "config.json"))

	cfg, err := repo.Get()
	if err != nil {
		t.Fatalf("get default: %v", err)
	}
	if cfg.WorkMinutes != 25 || cfg.ShortBreakMinutes != 5 {
		t.Fatalf("unexpected default: %+v", cfg)
	}
}

func TestConfigRepo_SaveAndGet(t *testing.T) {
	dir := t.TempDir()
	repo := NewConfigRepo(filepath.Join(dir, "config.json"))

	want := model.Config{WorkMinutes: 50, ShortBreakMinutes: 10, LongBreakMinutes: 20, LongBreakInterval: 3}
	if err := repo.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := repo.Get()
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestConfigRepo_AtomicWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	repo := NewConfigRepo(path)

	_ = repo.Save(model.Config{WorkMinutes: 25, ShortBreakMinutes: 5, LongBreakMinutes: 15, LongBreakInterval: 4})
	// tmp file should not linger
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("tmp file should not exist after save")
	}
}

func TestConfigRepo_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	repo := NewConfigRepo(filepath.Join(dir, "config.json"))

	if err := repo.Save(model.Config{WorkMinutes: 0, ShortBreakMinutes: 5, LongBreakMinutes: 15, LongBreakInterval: 4}); err == nil {
		t.Fatal("expected error for zero WorkMinutes")
	}
}
