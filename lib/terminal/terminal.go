// Terminal controls

package terminal

import (
	"fmt"
	"os"
	"strings"

	"github.com/mikeyQwn/doro/lib/ansi"
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

// Returns current terminal's width and height
func GetDimensions() (width int, height int, err error) {
	return term.GetSize(stdoutFd())
}

// Returns ANSI escape that moves cursor `n` lines up
func Up(n uint) string {
	return fmt.Sprintf("\033[%dA", n)
}

// Returns ANSI escape that moves cursor `n` lines down
func Down(n uint) string {
	return fmt.Sprintf("\033[%dB", n)
}

// Returns ANSI escape that moves cursor `n` lines down using linefeeds
func DownLF(n uint) string {
	return strings.Repeat(ansi.LINE_FEED, int(n))
}

func stdoutFd() int {
	return int(os.Stdout.Fd())
}
