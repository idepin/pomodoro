package model

// Config holds pomodoro timer settings.
type Config struct {
	WorkMinutes       int `json:"work_minutes"`
	ShortBreakMinutes int `json:"short_break_minutes"`
	LongBreakMinutes  int `json:"long_break_minutes"`
	// LongBreak triggers every N pomodoros.
	LongBreakInterval int `json:"long_break_interval"`
}

func DefaultConfig() Config {
	return Config{
		WorkMinutes:       25,
		ShortBreakMinutes: 5,
		LongBreakMinutes:  15,
		LongBreakInterval: 4,
	}
}
