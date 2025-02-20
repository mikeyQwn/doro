package ui

import (
	"strings"
	"unicode"

	"github.com/mikeyQwn/doro/lib/ansi"
	"github.com/mikeyQwn/doro/lib/terminal"
)

// Line formatting utilities
type Formatter struct {
	w int
}

// Creates a new formatter for the current terminal
// The formatter should be recreated after the terminal is resized
func NewFormatter() (*Formatter, error) {
	w, _, err := terminal.GetDimensions()
	if err != nil {
		return nil, err
	}

	return &Formatter{
		w,
	}, nil
}

// Pads spaces to the string to fit it in the center of the screen if it is able
// Returns the initial string if can't get the size of the terminal
func (f *Formatter) C(s string) string {
	printableCount := countPrintable(s)
	paddingLen := (f.w - printableCount) / 2
	if paddingLen <= 0 {
		return s
	}

	return strings.Repeat(" ", paddingLen) + s
}

// Wraps a string in ANSI escape that makes it bold and resets afterwards
func (f *Formatter) B(s string) string {
	return B(s)
}

// Wraps a string in ANSI escape that makes it bold and resets afterwards
func B(s string) string {
	return ansi.BOLD + s + ansi.RESET
}

func countPrintable(s string) int {
	cnt := 0
	isEscape := false
	for _, c := range s {
		if isEscape {
			if c == 'm' {
				isEscape = false
			}
			continue
		}

		if unicode.IsPrint(c) {
			cnt += 1
			continue
		}

		if c == '\033' {
			isEscape = true
		}
	}

	return cnt
}
