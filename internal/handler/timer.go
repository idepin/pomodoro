package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/idepin/pomodoro/internal/service"
)

type TimerHandler struct {
	svc *service.TimerService
}

func NewTimerHandler(svc *service.TimerService) *TimerHandler { return &TimerHandler{svc: svc} }

func (h *TimerHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/pomodoro", h.status)
	mux.HandleFunc("POST /api/pomodoro/start", h.start)
	mux.HandleFunc("POST /api/pomodoro/pause", h.pause)
	mux.HandleFunc("POST /api/pomodoro/resume", h.resume)
	mux.HandleFunc("POST /api/pomodoro/stop", h.stop)
}

func (h *TimerHandler) status(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.Status())
}

func (h *TimerHandler) start(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Mode string `json:"mode"`
	}
	if r.Body != nil && r.Body != http.NoBody {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
	}
	if err := h.svc.Start(in.Mode); err != nil {
		h.writeTimerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h.svc.Status())
}

func (h *TimerHandler) pause(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Pause(); err != nil {
		h.writeTimerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h.svc.Status())
}

func (h *TimerHandler) resume(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Resume(); err != nil {
		h.writeTimerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h.svc.Status())
}

func (h *TimerHandler) stop(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Stop(); err != nil {
		h.writeTimerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h.svc.Status())
}

func (h *TimerHandler) writeTimerError(w http.ResponseWriter, err error) {
	status := http.StatusConflict
	if errors.Is(err, service.ErrNotRunning) || errors.Is(err, service.ErrAlreadyStopped) {
		status = http.StatusBadRequest
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
