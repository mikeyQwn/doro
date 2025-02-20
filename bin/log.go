package bin

import (
	"os"

	"github.com/mikeyQwn/doro/lib/ansi"
)

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

func ExitSuccess() {
	os.Exit(0)
}
