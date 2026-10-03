package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/idepin/pomodoro/internal/model"
	"github.com/idepin/pomodoro/internal/repository"
	"github.com/idepin/pomodoro/internal/service"
)

type ConfigHandler struct {
	svc *service.ConfigService
}

func NewConfigHandler(svc *service.ConfigService) *ConfigHandler {
	return &ConfigHandler{svc: svc}
}

func (h *ConfigHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/config", h.get)
	mux.HandleFunc("POST /api/config", h.save)
}

func (h *ConfigHandler) get(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.svc.Get()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (h *ConfigHandler) save(w http.ResponseWriter, r *http.Request) {
	var cfg model.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := h.svc.Save(cfg); err != nil {
		if errors.Is(err, repository.ErrInvalidConfig) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	cfg, _ = h.svc.Get()
	writeJSON(w, http.StatusOK, cfg)
}
