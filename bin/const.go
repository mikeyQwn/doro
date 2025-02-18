package bin

import (
	"fmt"
	"strings"
	"time"

	"github.com/mikeyQwn/doro/lib"
	"github.com/mikeyQwn/doro/lib/terminal"
)

const (
	title                  = " Doro the pomodoro timer "
	focusedDurationLabel   = "Select focused work duration"
	breakDurationLabel     = "Select break duration"
	longBreakDurationLabel = "Select long break duration"
)

var (
	titleHorizontalBorder = fmt.Sprintf("+%s+", strings.Repeat("-", len(title)))
	initMsg               = boldCenteredCRLF.FormatLines(
		titleHorizontalBorder,
		"|"+title+"|",
		titleHorizontalBorder,
	)
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
	focusedDurationOffs     = 8
	focusedDurationSelector = lib.NewSelector(durations, focusedDurationOffs)

	breakDurationOffs     = 4
	breakDurationSelector = lib.NewSelector(durations, breakDurationOffs)

	longBreakDurationOffs     = 9
	longBreakDurationSelector = lib.NewSelector(durations, longBreakDurationOffs)
)

var centered = terminal.NewFormatBuilder().
	Center()

var centeredCRLF = centered.
	CRLF()

var boldCenteredCRLF = centeredCRLF.Bold()
