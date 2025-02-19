package bin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mikeyQwn/doro/lib"
	"github.com/mikeyQwn/doro/lib/terminal"
	"github.com/mikeyQwn/doro/lib/ui"
)

// Runs the pomodoro application
// Tranforms terminal into raw mode, executes the app
// Returns the terminal to it's original state before returning
func Run() {
	keyStreamCtx := context.Background()

	// Transform into raw mode and make sure to restore original state
	restore := Unwrap(terminal.IntoRaw())
	defer restore()
	AddOnExit(func() { restore() })

	ks := terminal.
		StdinIntoStream(keyStreamCtx, keyStreamBuffsize).
		HandleCtrlC(Exit)

	s := NewAppState(ks)
	stages := []func(){
		s.PrintInitMsg,
		s.FillMissingConfigFields,
		s.WaitForSpaceToStart,
		s.RunPomodoros,
	}

	for _, stageFn := range stages {
		stageFn()
		CheckErr(s.w.Done())
	}
}

func (s *AppState) PrintInitMsg() {
	s.w.WriteFmt(initMsg, fmtCRLF)
}

func (s *AppState) FillMissingConfigFields() {
	for _, selector := range s.cfg.GetMissingSelectors() {
		s.w.WriteFmt(selector.Label, fmtCenteredCRLF)
		*selector.ConfigRef = Unwrap(
			ui.RunSelector(selector.Selector, s.ks, lib.FormatDurationMinPrec),
		)
		s.w.WriteString(terminal.Down(2)).
			WriteString(terminal.ESCAPE_ERASE_RETURN)
	}

}

func (s *AppState) WaitForSpaceToStart() {
	s.w.WriteFmt(pressSpaceToStartMsg, fmtCenteredCRLF)
	s.ks.WaitKey(terminal.KEY_SPACE)
}

func (s *AppState) RunPomodoros() {
	s.w.WriteFmt(horizontalLineMsg, fmtCenteredCRLF)

	for i := 1; ; i++ {
		s.w.WriteFmt(fmt.Sprintf(pomodoroMsgTemplate, i), fmtCenteredCRLF)

		if spawnTimer(s.ks, focusedWorkLabel, s.cfg.focusedWorkDuration) {
			return
		}

		s.w.WriteString(terminal.Down(2)).
			WriteString(terminal.ESCAPE_ERASE_RETURN)

		if (i%4 != 0) && spawnTimer(s.ks, shortBreakLabel, s.cfg.breakDuration) {
			return
		}

		if (i%4 == 0) && spawnTimer(s.ks, longBreakLabel, s.cfg.longBreakDuration) {
			return
		}

		s.w.WriteString(terminal.Down(2)).
			WriteString(terminal.ESCAPE_ERASE_RETURN)
	}
}

func spawnTimer(keyStream <-chan terminal.Key, label string, total time.Duration) bool {
	fmt.Print(strings.Repeat(terminal.ESCAPE_CR_LF, 2))
	fmt.Println(fmtCenteredCRLF.Format("Press " + terminal.Bold("[space]") + " to pause"))
	fmt.Print(terminal.Up(4))

	start := time.Now()
	addSpare := false
	fmt.Print(formatProgress(label, total, start, addSpare))

	for {
		select {
		case k := <-keyStream:
			_ = k
		case <-time.After(time.Second * 1):
			if time.Now().Sub(start) > total {
				fmt.Print(formatProgress(label, total, start, addSpare))
				return false
			}

			addSpare = !addSpare
			fmt.Print(formatProgress(label, total, start, addSpare))
		}
	}
}

func formatProgress(label string, total time.Duration, start time.Time, addSpare bool) string {
	totalMin := total.Minutes()
	elapsedMin := time.Now().Sub(start).Minutes()
	msg := fmt.Sprintf("%s: %s %d/%dm\r",
		label,
		formatProgressBar(elapsedMin/totalMin, 25, addSpare),
		int(elapsedMin),
		int(totalMin))
	return fmtCentered.Format(msg)
}

func formatProgressBar(completion float64, width uint, addSpare bool) string {
	completedLen := uint(float64(width) * completion)
	spare := ""
	missingLen := width - completedLen
	if addSpare && completedLen < width {
		spare = "-"
		missingLen -= 1
	}
	completed := strings.Repeat("=", int(completedLen))
	missing := strings.Repeat(" ", int(missingLen))
	return "[" + completed + spare + missing + "]"
}
