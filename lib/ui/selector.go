// High-level abstractions over terminal in raw mode

package ui

import (
	"fmt"

	"github.com/mikeyQwn/doro/lib"
	"github.com/mikeyQwn/doro/lib/terminal"
)

// Runs the selector on the current line
// Transforms value inside the selector with provided `labelFunc`
func RunSelector[T any](s *lib.Selector[T], labelFunc func(T) string) (selectedOption T, shouldExit bool, e error) {
	for {
		fmt.Print(terminal.Center(fmt.Sprintf("< %s >", labelFunc(s.Curr()))))
		k, err := terminal.CaptupreKey()
		if err != nil {
			var empty T
			return empty, false, err
		}

		switch k {
		case terminal.KEY_ARROW_LEFT:
			_ = s.Prev()
		case terminal.KEY_ARROW_RIGHT:
			_ = s.Next()
		case terminal.KEY_ENTER:
			return s.Curr(), false, nil
		case terminal.KEY_CTRL_C:
			var t T
			return t, true, nil
		}
		fmt.Print(terminal.ESCAPE_ERASE_LINE + terminal.ESCAPE_CARRIAGE_RETURN)
	}
}

// Waits until user presses the given `key` or Ctrl-C
// Returns if Ctrl-C was pressed
func WaitKey(key terminal.Key) (bool, error) {
	for {
		k, err := terminal.CaptupreKey()
		if err != nil {
			return false, err
		}

		if k == terminal.KEY_CTRL_C {
			return true, nil
		}

		if k == key {
			return false, nil
		}
	}
}
