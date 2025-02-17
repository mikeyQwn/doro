package terminal

import (
	"strings"
	"unicode"
)

const (
	ESCAPE_RESET = "\033[1m"
	ESCAPE_BOLD  = "\033[1m"
)

// Wraps a string in ansii escape that makes it Bold and resets afterwards
func Bold(s string) string {
	return ESCAPE_BOLD + s + ESCAPE_RESET
}

// Pads spaces to the string to fit it in the center of the screen if it is able
// Returns the initial string if can't get the size of the terminal
func Center(s string) string {
	w, _, err := GetTerminalDimensions()
	if err != nil {
		return s
	}
	printableCount := countPrintable(s)
	paddingLen := (w - printableCount) / 2
	return strings.Repeat(" ", paddingLen) + s
}

// Adds \r\n to the end of the string
func CRLF(s string) string {
	return s + "\r\n"
}

func countPrintable(s string) int {
	cnt := 0
	for _, c := range s {
		if unicode.IsPrint(c) {
			cnt += 1
		}
	}

	return cnt
}

// Formats texts with selected options. By default, everything is deselected
type FormatBuilder struct {
	center bool
	bold   bool
	crlf   bool
}

// Instantiates `FormatBuilder`
// As it's size is not that big, copying it in every call is preferable
func NewFormatBuilder() FormatBuilder {
	return FormatBuilder{}
}

// Makes the otput of the returned formatter bold
func (f FormatBuilder) Bold() FormatBuilder {
	f.bold = true
	return f
}

// Makes the otput of the returned formatter centered
func (f FormatBuilder) Center() FormatBuilder {
	f.center = true
	return f
}

// Makes the otput of the returned formatter escaped by \r\n
func (f FormatBuilder) CRLF() FormatBuilder {
	f.crlf = true
	return f
}

// Returs a string, formatted with selected options
func (f FormatBuilder) Format(s string) string {
	if f.center {
		s = Center(s)
	}
	if f.bold {
		s = Bold(s)
	}
	if f.crlf {
		s = CRLF(s)
	}
	return s
}

// Apllies the formatter to each individual line separately
func (f FormatBuilder) FormatLines(lines ...string) string {
	sb := strings.Builder{}

	for _, line := range lines {
		_, _ = sb.WriteString(f.Format(line))
	}

	return sb.String()
}
