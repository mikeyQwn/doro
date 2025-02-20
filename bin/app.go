package bin

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mikeyQwn/doro/lib"
	input "github.com/mikeyQwn/doro/lib/input"
	tm "github.com/mikeyQwn/doro/lib/terminal"

	"github.com/mikeyQwn/doro/lib/ui"
)

// Runs the pomodoro application
// Tranforms terminal into raw mode, executes the app
// Returns the terminal to it's original state before returning
func Run() {
	keyStreamCtx := context.Background()

	// Transform into raw mode and make sure to restore original state
	restore := Unwrap(tm.IntoRaw())
	defer restore()
	AddOnExit(func() { restore() })

	ks := input.
		StdinIntoStream(keyStreamCtx, keyStreamBuffsize).
		HandleCtrlC(Exit)

	s := NewAppState(ks)

	CheckErr(s.InitMsg().Run())
	Unwrap(io.WriteString(s.wr, tm.DownLF(1)))

	for _, widget := range s.ConfigSelectors() {
		CheckErr(widget.Run())
		Unwrap(io.WriteString(s.wr, tm.DownLF(1)))
	}

	CheckErr(s.WaitForSpace().Run())
	Unwrap(io.WriteString(s.wr, tm.DownLF(1)))

	for n := 1; ; n++ {
		widget := s.CreatePomodoro(n)
		CheckErr(widget.Run())
		<-time.After(time.Second)
		Unwrap(io.WriteString(s.wr, tm.DownLF(1)))
	}
}

func (s *AppState) InitMsg() *ui.Widget {
	return ui.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		return []string{
			f.C(tm.B(titleHorizontalBorder)),
			f.C(tm.B("|" + title + "|")),
			f.C(tm.B(titleHorizontalBorder)),
		}, true
	}).WithWriter(s.wr)
}

func (s *AppState) ConfigSelectors() []*ui.Widget {
	selectors := s.cfg.GetMissingSelectors()
	widgets := make([]*ui.Widget, 0, len(selectors))

	for _, selector := range selectors {
		w := s.ConfigSelector(selector)
		widgets = append(widgets, w)
	}

	return widgets
}

func (s *AppState) ConfigSelector(selector ConfigSelector) *ui.Widget {
	sel := selector.Selector
	done := false

	return s.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		val := sel.Curr()

		if done {
			*selector.ConfigRef = val
		}

		msg := fmt.Sprintf("< %s >", lib.FormatDurationMinPrec(val))

		return []string{
			f.C(selector.Label),
			f.C(msg),
		}, done
	}).AddKeyHandler(func(k input.Key) { sel.Prev() }, input.KEY_ARROW_LEFT).
		AddKeyHandler(func(k input.Key) { sel.Next() }, input.KEY_ARROW_RIGHT).
		AddKeyHandler(func(k input.Key) { done = true }, input.KEY_ENTER)
}

func (s *AppState) WaitForSpace() *ui.Widget {
	done := false
	return s.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		return []string{
			f.C(pressSpaceToStartMsg),
		}, done
	}).AddKeyHandler(func(k input.Key) { done = true }, input.KEY_SPACE)
}

func (s *AppState) CreatePomodoro(n int) *ui.Widget {
	done := false
	addDot := false
	pomodoroMsg := fmt.Sprintf(pomodoroMsgTemplate, n)
	w := 16

	labelA := focusedWorkLabel
	completionA := formatProgressBar(0.0, uint(w), false)
	durationA := s.cfg.breakDuration
	timerA := lib.NewTimer(durationA)

	labelB := shortBreakLabel
	completionB := formatProgressBar(0.0, uint(w), false)
	durationB := s.cfg.breakDuration
	if n%4 == 0 {
		labelB = longBreakLabel
		durationB = s.cfg.longBreakDuration
	}
	timerB := lib.NewPaused(durationB)
	activeTimer := timerA

	return s.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		return []string{
			f.C(pomodoroMsg),
			"",
			f.C(labelA + ": " + completionA + " " + lib.FormatPercent(timerA.Progress())),
			"",
			f.C(labelB + ": " + completionB + " " + lib.FormatPercent(timerB.Progress())),
			"",
			f.C("Press " + tm.B("[space]") + " to pause"),
		}, done
	}).AddTimedHandler(func() {
		if activeTimer.IsPaused() {
			return
		}

		addDot = !addDot
		progressA := timerA.Progress()
		if timerA.IsFinished() && timerB.IsPaused() && timerB.Elapsed() == 0 {
			timerB.Unpause()
			activeTimer = timerB
		}

		progressB := timerB.Progress()
		if timerB.IsFinished() {
			done = true
			return
		}

		completionA = formatProgressBar(progressA, uint(w), !timerA.IsPaused() && addDot)
		completionB = formatProgressBar(progressB, uint(w), !timerB.IsPaused() && addDot)
	}, time.Second*1).AddKeyHandler(func(k input.Key) { activeTimer.Toggle() }, input.KEY_SPACE)
}

func formatProgressBar(completion float64, width uint, addDot bool) string {
	completedLen := uint(float64(width) * completion)
	spare := ""
	missingLen := width - completedLen
	if addDot && completedLen < width {
		spare = "-"
		missingLen -= 1
	}
	completed := strings.Repeat("=", int(completedLen))
	missing := strings.Repeat(" ", int(missingLen))
	return "[" + completed + spare + missing + "]"
}
