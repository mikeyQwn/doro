package bin

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mikeyQwn/doro/lib/ansi"
	input "github.com/mikeyQwn/doro/lib/input"
	tm "github.com/mikeyQwn/doro/lib/terminal"

	"github.com/mikeyQwn/doro/lib/ui"
)

// Runs the pomodoro application
// Tranforms terminal into raw mode, executes the app
// Returns the terminal to it's original state before returning
func Run() error {
	cfg := ParseConfig()
	if cfg.printVersion {
		fmt.Println(os.Args[0], currentVersion)
		return nil
	}

	keyStreamCtx, keyStreamCtxCancel := context.WithCancel(context.Background())
	defer keyStreamCtxCancel()

	// Transform into raw mode and make sure to restore original state
	restore, err := tm.IntoRaw()
	if err != nil {
		return err
	}
	defer func() { _ = restore() }()

	ks := input.
		StdinIntoStream(keyStreamCtx, keyStreamBuffsize).
		HandleCtrlC(func() { _ = restore(); os.Exit(0) })

	s := NewAppState(ks, cfg)

	widgets := [][]*ui.Widget{
		{s.InitMsg()},
		s.ConfigSelectors(),
		{s.WaitForSpace()},
	}

	// Execute initial setup and wait for space press
	for _, widgetRow := range widgets {
		for _, widget := range widgetRow {
			if err := widget.Run(); err != nil {
				return err
			}

			if _, err := io.WriteString(s.wr, tm.DownLF(1)); err != nil {
				return err
			}
		}
	}

	// Run the pomodoro loop
	for n := 1; ; n++ {
		widget := s.CreatePomodoro(n)
		if err := widget.Run(); err != nil {
			return err
		}

		<-time.After(time.Second)
		if _, err := io.WriteString(s.wr, tm.Up(3)+ansi.ERASE_LINE); err != nil {
			return err
		}
	}
}

// Prints the main title
func (s *AppState) InitMsg() *ui.Widget {
	return s.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		return []string{
			f.C(f.B(titleHorizontalBorder)),
			f.C(f.B("| " + title + " |")),
			f.C(f.B(titleHorizontalBorder)),
		}, true
	})
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

// Selects a value for a particular config field
func (s *AppState) ConfigSelector(selector ConfigSelector) *ui.Widget {
	sel := selector.Selector
	done := false

	return s.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		val := FormatDurationMinPrec(sel.Curr())
		return []string{
			f.C(selector.Label),
			f.C("< " + f.B(val) + " >"),
		}, done
	}).AddKeyHandler(func(k input.Key) { sel.Prev() }, input.KEY_ARROW_LEFT).
		AddKeyHandler(func(k input.Key) { sel.Next() }, input.KEY_ARROW_RIGHT).
		AddKeyHandler(func(k input.Key) { *selector.ConfigRef = sel.Curr(); done = true }, input.KEY_ENTER)
}

// Prints a hint and waits until space is pressed
func (s *AppState) WaitForSpace() *ui.Widget {
	done := false
	return s.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		return []string{
			f.C("Press " + f.B("[space]") + " to start!"),
		}, done
	}).AddKeyHandler(func(k input.Key) { done = true }, input.KEY_SPACE)
}

func (s *AppState) CreatePomodoro(n int) *ui.Widget {
	addDot := false

	isLong := n%4 == 0
	pd := NewPomodoro(s.cfg, isLong)

	pomodoroMsg := fmt.Sprintf(pomodoroMsgTemplate, n)

	workCompletion := formatProgressBar(0.0, progressBarWidth, false)
	breakCompletion := formatProgressBar(0.0, progressBarWidth, false)

	return s.NewWidget(func(f *ui.Formatter) ([]string, bool) {
		workCompletion = formatProgressBar(pd.WorkProgress(), progressBarWidth, pd.WorkRunning() && addDot)
		breakCompletion = formatProgressBar(pd.BreakProgress(), progressBarWidth, pd.BreakRunning() && addDot)

		workTimer := FormatTimer(pd.WorkElapsed(), pd.WorkDuration())
		breakTimer := FormatTimer(pd.BreakElapsed(), pd.BreakDuration())

		title := pomodoroMsg

		isPaused := pd.IsPaused()
		pauseMsg := "Press " + f.B("[space]") + " "
		if isPaused {
			title += " " + f.B("[PAUSED]")
			pauseMsg += "to unpause"
		} else {
			pauseMsg += "to pause"
		}

		skipMsg := "Press " + f.B("[s]") + " to skip the current session"

		return []string{
			f.C(title),
			"",
			f.C(fmt.Sprintf("%s: %s %s (%s)", pd.WorkLabel(), workCompletion, workTimer, FormatPercent(pd.WorkProgress()))),
			"",
			f.C(fmt.Sprintf("%s: %s %s (%s)", pd.BreakLabel(), breakCompletion, breakTimer, FormatPercent(pd.BreakProgress()))),
			"",
			f.C(pauseMsg),
			"",
			f.C(skipMsg),
		}, pd.IsFinished()
	}).AddTimedHandler(func() {
		if pd.IsPaused() {
			return
		}
		addDot = !addDot
		pd.Update()

	}, time.Second*1).
		AddKeyHandler(func(k input.Key) { pd.TogglePause() }, input.KEY_SPACE).
		AddKeyHandler(func(k input.Key) { pd.TogglePause(); pd.NextTask() }, input.KEY_S)
}

func formatProgressBar(completion float64, width uint, addDot bool) string {
	bar := make([]byte, width+2)
	bar[0] = '['
	bar[len(bar)-1] = ']'
	completedLen := uint(float64(width) * completion)
	for i := range width {
		if i < completedLen {
			bar[1+i] = '='
		} else {
			bar[1+i] = ' '
		}
	}
	if addDot && completedLen < width {
		bar[1+completedLen] = '-'
	}

	return string(bar)
}
