package bin

import (
	"os"

	"github.com/mikeyQwn/doro/lib/ansi"
)

// Prints the message to `stderr` and exits with libc's EXIT_FAILURE
func Fatal(msg string) {
	_, _ = os.Stderr.WriteString(
		ansi.RED +
			"Fatal error:" +
			ansi.RESET +
			" " +
			msg,
	)

	os.Exit(1)
}

// Exits with libc's EXIT_SUCCESS status
func ExitSuccess() {
	os.Exit(0)
}
