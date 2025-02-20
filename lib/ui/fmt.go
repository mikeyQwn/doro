package ui

import (
	"strings"
	"unicode"

	"github.com/mikeyQwn/doro/lib/terminal"
)

// Some formatting requires state that is persistent in a single
// Update call, so this object eliviates unnecessary calls to
// retreive that state
type Formatter struct {
	w, h int
}

func NewFormatter() (*Formatter, error) {
	w, h, err := terminal.GetDimensions()
	if err != nil {
		return nil, err
	}

	return &Formatter{
		w, h,
	}, nil
}

// Does the same thing as as terminal.C
func (f *Formatter) C(s string) string {
	printableCount := countPrintable(s)
	paddingLen := (f.w - printableCount) / 2
	if paddingLen <= 0 {
		return s
	}

	return strings.Repeat(" ", paddingLen) + s
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
