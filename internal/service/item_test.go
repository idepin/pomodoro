package service

import (
	"testing"

	"github.com/idepin/pomodoro/internal/repository"
)

func TestItemCRUD(t *testing.T) {
	svc := NewItemService(repository.NewMemoryItemRepo())

	if _, err := svc.Create("1", "sword"); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := svc.Get("1")
	if err != nil || got.Name != "sword" {
		t.Fatalf("get = %+v, err = %v", got, err)
	}
	if len(svc.List()) != 1 {
		t.Fatalf("list len = %d, want 1", len(svc.List()))
	}
	if err := svc.Delete("1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.Get("1"); err == nil {
		t.Fatal("expected not found after delete")
	}
	if _, err := svc.Create("", "x"); err == nil {
		t.Fatal("expected validation error")
	}
}
