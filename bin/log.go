package bin

import (
	"os"
	"sync"

	"github.com/mikeyQwn/doro/lib/terminal"
)

var onExitLock sync.Mutex
var onExit = []func(){}

func AddOnExit(f func()) {
	onExitLock.Lock()
	defer onExitLock.Unlock()

	onExit = append(onExit, f)
}

func CheckErr(err error) {
	if err != nil {
		Fatal(err.Error())
	}
}

func Unwrap[T any](a T, err error) T {
	if err != nil {
		Fatal(err.Error())
	}

	return a
}

func Unwrap2[T any, U any](a T, b U, err error) (T, U) {
	if err != nil {
		Fatal(err.Error())
	}

	return a, b
}

func Fatal(msg string) {
	_, _ = os.Stderr.WriteString(
		terminal.ESCAPE_RED +
			"Fatal error:" +
			terminal.ESCAPE_RESET +
			" " +
			msg,
	)

	ExitWith(1)
}

func ExitWith(status int) {
	onExitLock.Lock()
	defer onExitLock.Unlock()
	for _, f := range onExit {
		f()
	}

	os.Exit(status)
}

func Exit() {
	ExitWith(0)
}
