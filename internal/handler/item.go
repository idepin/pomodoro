package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/idepin/pomodoro/internal/repository"
	"github.com/idepin/pomodoro/internal/service"
)

type ItemHandler struct {
	svc *service.ItemService
}

func NewItemHandler(svc *service.ItemService) *ItemHandler {
	return &ItemHandler{svc: svc}
}

// Register wires REST routes. Go 1.22+ ServeMux supports "METHOD /path".
func (h *ItemHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/items", h.list)
	mux.HandleFunc("GET /api/items/{id}", h.get)
	mux.HandleFunc("POST /api/items", h.create)
	mux.HandleFunc("DELETE /api/items/{id}", h.delete)
}

func (h *ItemHandler) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.List())
}

func (h *ItemHandler) get(w http.ResponseWriter, r *http.Request) {
	item, err := h.svc.Get(r.PathValue("id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	item, err := h.svc.Create(in.ID, in.Name)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.PathValue("id")); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeErr(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
