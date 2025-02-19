package terminal

import (
	"os"

	"golang.org/x/term"
)

type TerminalRestoreFunc func() error

// Puts the terminal into raw mode and returns a function that restores
// it to normal state. The returned function almost always should be called
// before program exits
func IntoRaw() (TerminalRestoreFunc, error) {
	s, err := term.MakeRaw(stdoutFd())
	if err != nil {
		return nil, err
	}

	restoreFn := func() error {
		return term.Restore(stdoutFd(), s)
	}

	return restoreFn, nil
}

func GetDimensions() (width int, height int, err error) {
	return term.GetSize(stdoutFd())
}

func stdoutFd() int {
	return int(os.Stdout.Fd())
}
