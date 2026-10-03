package service

import (
	"errors"
	"sync"
	"time"

	"github.com/idepin/pomodoro/internal/model"
)

var (
	ErrAlreadyRunning = errors.New("timer already running")
	ErrNotRunning     = errors.New("timer not running")
	ErrAlreadyStopped = errors.New("timer already stopped")
)

type TimerService struct {
	mu        sync.Mutex
	state     model.TimerState
	mode      string
	remaining int       // seconds left
	startedAt time.Time // when last resumed
	configSvc *ConfigService
}

func NewTimerService(configSvc *ConfigService) *TimerService {
	return &TimerService{
		state:     model.TimerStopped,
		configSvc: configSvc,
	}
}

// elapsed returns seconds elapsed since last resume (only valid when running).
func (s *TimerService) elapsed() int {
	return int(time.Since(s.startedAt).Seconds())
}

func (s *TimerService) Status() model.Timer {
	s.mu.Lock()
	defer s.mu.Unlock()

	remaining := s.remaining
	if s.state == model.TimerRunning {
		remaining -= s.elapsed()
		if remaining < 0 {
			remaining = 0
		}
	}
	var startedAt *time.Time
	if s.state == model.TimerRunning {
		t := s.startedAt
		startedAt = &t
	}
	return model.Timer{
		State:     s.state,
		Mode:      s.mode,
		Remaining: remaining,
		StartedAt: startedAt,
	}
}

// Start begins a new pomodoro. mode: "work" | "short_break" | "long_break"
func (s *TimerService) Start(mode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == model.TimerRunning {
		return ErrAlreadyRunning
	}

	cfg, err := s.configSvc.Get()
	if err != nil {
		return err
	}

	var duration int
	switch mode {
	case "short_break":
		duration = cfg.ShortBreakMinutes * 60
	case "long_break":
		duration = cfg.LongBreakMinutes * 60
	default:
		mode = "work"
		duration = cfg.WorkMinutes * 60
	}

	s.mode = mode
	s.remaining = duration
	s.state = model.TimerRunning
	s.startedAt = time.Now()
	return nil
}

func (s *TimerService) Pause() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != model.TimerRunning {
		return ErrNotRunning
	}
	s.remaining -= s.elapsed()
	if s.remaining < 0 {
		s.remaining = 0
	}
	s.state = model.TimerPaused
	return nil
}

func (s *TimerService) Resume() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != model.TimerPaused {
		return ErrNotRunning
	}
	s.startedAt = time.Now()
	s.state = model.TimerRunning
	return nil
}

func (s *TimerService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == model.TimerStopped {
		return ErrAlreadyStopped
	}
	s.state = model.TimerStopped
	s.remaining = 0
	s.mode = ""
	return nil
}
