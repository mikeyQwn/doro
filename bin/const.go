package bin

import (
	"fmt"
	"strings"
	"time"

	"github.com/mikeyQwn/doro/lib"
)

const (
	currentVersion         = "v0.0.1"
	title                  = "Doro the pomodoro timer"
	focusedDurationLabel   = "Select focused work duration"
	focusedWorkLabel       = "Focused work"
	breakDurationLabel     = "Select break duration"
	shortBreakLabel        = "Short break"
	longBreakDurationLabel = "Select long break duration"
	longBreakLabel         = "Long break"

	pomodoroMsgTemplate = "Pomodoro %d"
)

var (
	titleHorizontalBorder = fmt.Sprintf("+%s+", strings.Repeat("-", len(title)+2))
)

var durations = []time.Duration{
	time.Minute * 1,
	time.Minute * 2,
	time.Minute * 3,
	time.Minute * 4,
	time.Minute * 5,
	time.Minute * 10,
	time.Minute * 15,
	time.Minute * 20,
	time.Minute * 25,
	time.Minute * 30,
	time.Minute * 45,
	time.Minute * 60,
	time.Minute * 120,
}

var (
	focusedWorkDurationOffs     = 8
	focusedWorkDurationSelector = lib.NewSelector(durations, focusedWorkDurationOffs)

	shortBreakDurationOffs     = 4
	shortBreakDurationSelector = lib.NewSelector(durations, shortBreakDurationOffs)

	longBreakDurationOffs     = 9
	longBreakDurationSelector = lib.NewSelector(durations, longBreakDurationOffs)
)

const (
	keyStreamBuffsize = 16
	progressBarWidth  = uint(16)
)
