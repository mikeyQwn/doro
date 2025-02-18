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
	breakDuration       time.Duration
	longBreakDuration   time.Duration
}

func NewConfigSelectors(c *Config) []ConfigSelector {
	return []ConfigSelector{
		ConfigSelector{
			Label:     focusedDurationLabel,
			Selector:  focusedDurationSelector,
			ConfigRef: &c.focusedWorkDuration,
		},
		ConfigSelector{
			Label:     breakDurationLabel,
			Selector:  breakDurationSelector,
			ConfigRef: &c.breakDuration,
		},
		ConfigSelector{
			Label:     longBreakDurationLabel,
			Selector:  longBreakDurationSelector,
			ConfigRef: &c.longBreakDuration,
		},
	}
}
