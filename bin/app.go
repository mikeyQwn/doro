package bin

import (
	"fmt"
	"log"
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

	fmt.Println(centeredCRLF.Format("--------------"))
	for i := 1; ; i++ {
		fmt.Println(centeredCRLF.Format(fmt.Sprintf("Pomodoro %d", i)))
		fmt.Print(terminal.Down(2))
		fmt.Println(centeredCRLF.Format("Press " + terminal.Bold("[space]") + " to pause"))
		fmt.Print(terminal.Up(4))
		fmt.Print(centered.Format(fmt.Sprintf("Focused work: [=========== ] 0/25m\r")))
		<-time.After(time.Second * 3)
		fmt.Print(centered.Format(fmt.Sprintf("Focused work: [=========== ] 1/25m\r")))
		time.Sleep(time.Hour)
	}
}
