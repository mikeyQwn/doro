package bin

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mikeyQwn/doro/lib/terminal"
	"github.com/mikeyQwn/doro/lib/ui"
)

func durationToMinuteStirng(d time.Duration) string {
	return fmt.Sprintf("%dm", int(d.Minutes()))
}

func unwrap[T any](res T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func unwrap2[T any, U any](a T, b U, err error) (T, U) {
	if err != nil {
		log.Fatal(err)
	}

	return a, b
}

func Run() {
	c := Config{}
	selectors := NewConfigSelectors(&c)

	fmt.Print(initMsg)

	fmt.Print(terminal.ESCAPE_CR_LF)
	for _, s := range selectors {
		var shouldExit bool

		fmt.Print(centeredCRLF.Format(s.Label))
		*s.ConfigRef, shouldExit = unwrap2(ui.RunSelector(s.Selector, durationToMinuteStirng))
		if shouldExit {
			return
		}

		fmt.Println(terminal.ESCAPE_CR_LF)

	}

	fmt.Println(centeredCRLF.Format("Press " + terminal.Bold("[space]") + " to start!"))
	if unwrap(ui.WaitKey(terminal.KEY_SPACE)) {
		return
	}

	keyStream := terminal.StdinIntoStream(context.Background(), 16)

	fmt.Println(centeredCRLF.Format("--------------"))
	for i := 1; ; i++ {
		fmt.Println(centeredCRLF.Format(fmt.Sprintf("Pomodoro %d", i)))
		if spawnTimer(keyStream, "Focused work", c.focusedWorkDuration) {
			return
		}
		fmt.Print(terminal.ESCAPE_CARRIAGE_RETURN + terminal.ESCAPE_CR_LF + terminal.ESCAPE_CR_LF + terminal.ESCAPE_ERASE_LINE)

		if (i%4 != 0) && spawnTimer(keyStream, "Short break", c.breakDuration) {
			return
		}
		if (i%4 == 0) && spawnTimer(keyStream, "Long break", c.longBreakDuration) {
			return
		}
		fmt.Print(strings.Repeat(terminal.ESCAPE_CR_LF, 2))
		fmt.Print(terminal.ESCAPE_ERASE_LINE + terminal.ESCAPE_CARRIAGE_RETURN)
	}
}

func spawnTimer(keyStream <-chan terminal.Key, label string, total time.Duration) bool {
	fmt.Print(strings.Repeat(terminal.ESCAPE_CR_LF, 2))
	fmt.Println(centeredCRLF.Format("Press " + terminal.Bold("[space]") + " to pause"))
	fmt.Print(terminal.Up(4))

	start := time.Now()
	addSpare := false
	fmt.Print(formatProgress(label, total, start, addSpare))

	for {
		select {
		case k := <-keyStream:
			if k == terminal.KEY_CTRL_C {
				return true
			}
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
	return centered.Format(msg)
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
