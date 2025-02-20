package bin

import (
	"time"

	"github.com/mikeyQwn/doro/lib"
)

type ConfigSelector struct {
	Label     string
	Selector  *lib.Selector[time.Duration]
	ConfigRef *time.Duration
}

type Config struct {
	focusedWorkDuration time.Duration
	shortBreakDuration  time.Duration
	longBreakDuration   time.Duration
}

func (c *Config) GetMissingSelectors() []ConfigSelector {
	return []ConfigSelector{
		ConfigSelector{
			Label:     focusedDurationLabel,
			Selector:  focusedDurationSelector,
			ConfigRef: &c.focusedWorkDuration,
		},
		ConfigSelector{
			Label:     breakDurationLabel,
			Selector:  breakDurationSelector,
			ConfigRef: &c.shortBreakDuration,
		},
		ConfigSelector{
			Label:     longBreakDurationLabel,
			Selector:  longBreakDurationSelector,
			ConfigRef: &c.longBreakDuration,
		},
	}
}
