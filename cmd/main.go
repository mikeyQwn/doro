package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/mikeyQwn/doro/lib/terminal"
)

func main() {
	res, err := terminal.IntoRaw()
	if err != nil {
		log.Fatal(err)
	}
	defer res()

	boldCenteredCRLF := terminal.NewFormatBuilder().
		Bold().
		Center().
		CRLF()

	title := "Doro the pomodoro timer"

	titleHorizontalBorder := fmt.Sprintf("+%s+", strings.Repeat("-", len(title)))

	msg := boldCenteredCRLF.FormatLines(
		titleHorizontalBorder,
		"|"+title+"|",
		titleHorizontalBorder,
	)

	fmt.Print(msg)
}
