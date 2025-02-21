package bin

import (
	"flag"
	"os"
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
	printVersion        bool
	keepDefaults        bool
}

func ParseConfig() *Config {
	c := &Config{}
	s := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	s.DurationVar(&c.focusedWorkDuration, "f", 0, "specify focused work duration")
	s.DurationVar(&c.shortBreakDuration, "s", 0, "specify short break duration")
	s.DurationVar(&c.longBreakDuration, "l", 0, "specify long break duration")
	s.BoolVar(&c.keepDefaults, "d", false, "keep the default configuration")
	s.BoolVar(&c.printVersion, "v", false, "print version info and exit")
	_ = s.Parse(os.Args[1:])
	if c.keepDefaults {
		c.SetDefaults()
	}
	return c
}

func (c *Config) SetDefaults() {
	c.focusedWorkDuration = focusedWorkDurationSelector.Curr()
	c.shortBreakDuration = shortBreakDurationSelector.Curr()
	c.longBreakDuration = longBreakDurationSelector.Curr()
}

func (c *Config) GetMissingSelectors() []ConfigSelector {
	s := []ConfigSelector{}
	if c.keepDefaults {
		return s
	}

	if c.focusedWorkDuration == 0 {
		s = append(s,
			ConfigSelector{
				Label:     focusedDurationLabel,
				Selector:  focusedWorkDurationSelector,
				ConfigRef: &c.focusedWorkDuration,
			},
		)
	}
	if c.shortBreakDuration == 0 {
		s = append(s,
			ConfigSelector{
				Label:     breakDurationLabel,
				Selector:  shortBreakDurationSelector,
				ConfigRef: &c.shortBreakDuration,
			},
		)
	}

	if c.longBreakDuration == 0 {
		s = append(s,
			ConfigSelector{
				Label:     longBreakDurationLabel,
				Selector:  longBreakDurationSelector,
				ConfigRef: &c.longBreakDuration,
			},
		)
	}

	return s
}
