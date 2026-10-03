package model

import "time"

type TimerState string

const (
	TimerStopped TimerState = "stopped"
	TimerRunning TimerState = "running"
	TimerPaused  TimerState = "paused"
)

type Timer struct {
	State     TimerState `json:"state"`
	Mode      string     `json:"mode"`
	Remaining int        `json:"remaining_seconds"`
	StartedAt *time.Time `json:"started_at,omitempty"`
}
