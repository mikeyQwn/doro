// High-level abstractions over terminal in raw mode

package ui

import (
	"fmt"
	"os"

	"github.com/mikeyQwn/doro/lib"
	"github.com/mikeyQwn/doro/lib/terminal"
)

// Runs the selector on the current line
// Transforms value inside the selector with provided `labelFunc`
func RunSelector[T any](
	s *lib.Selector[T], keyStream terminal.KeyStream, labelFunc func(T) string,
) (T, error) {
	w := os.Stdout

	format := terminal.NewFormatBuilder().Center().Bold()

	for {
		msg := fmt.Sprintf("< %s >", labelFunc(s.Curr()))

		w.WriteString(format.Format(msg))

		k := <-keyStream
		switch k {
		case terminal.KEY_ARROW_LEFT:
			_ = s.Prev()
		case terminal.KEY_ARROW_RIGHT:
			_ = s.Next()
		case terminal.KEY_ENTER:
			return s.Curr(), nil
		}

		w.WriteString(terminal.ESCAPE_ERASE_RETURN)
	}
}
