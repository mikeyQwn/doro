package terminal

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	ESCAPE_RESET         = "\033[0m"
	ESCAPE_BOLD          = "\033[1m"
	ESCAPE_ERASE_LINE    = "\033[2K"
	ESCAPE_UP_LINE       = "\033[A"
	ESCAPE_UP_START_LINE = "\033[F"

	ESCAPE_RED = "\033[31m"

	ESCAPE_CARRIAGE_RETURN = "\r"
	ESCAPE_LINE_FEED       = "\n"

	ESCAPE_ERASE_RETURN = ESCAPE_ERASE_LINE + ESCAPE_CARRIAGE_RETURN
	ESCAPE_CR_LF        = ESCAPE_CARRIAGE_RETURN + ESCAPE_LINE_FEED
)

// Wraps a string in ANSI escape that makes it bold and resets afterwards
func B(s string) string {
	return ESCAPE_BOLD + s + ESCAPE_RESET
}

// Pads spaces to the string to fit it in the center of the screen if it is able
// Returns the initial string if can't get the size of the terminal
func C(s string) string {
	w, _, err := GetDimensions()
	if err != nil {
		return s
	}
	printableCount := countPrintable(s)
	paddingLen := (w - printableCount) / 2
	if paddingLen <= 0 {
		return s
	}

	return strings.Repeat(" ", paddingLen) + s
}

// Adds \r to the end of the string
func CR(s string) string {
	return s + ESCAPE_CARRIAGE_RETURN
}

// Adds \n to the end of the string
func LF(s string) string {
	return s + ESCAPE_LINE_FEED
}

// Adds \r\n to the end of the string
func CRLF(s string) string {
	return s + ESCAPE_CR_LF
}

// Returns ANSI escape that moves cursor `n` lines up
func Up(n uint) string {
	return fmt.Sprintf("\033[%dA", n)
}

// Returns ANSI escape that moves cursor `n` lines down
func Down(n uint) string {
	return fmt.Sprintf("\033[%dB", n)
}

// Returns ANSI escape that moves cursor `n` lines down using \n
func DownLF(n uint) string {
	return strings.Repeat(ESCAPE_LINE_FEED, int(n))
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

// Formats texts with selected options. By default, everything is deselected
type FormatBuilder struct {
	center bool
	bold   bool
	cr     bool
	lf     bool
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

// Makes the otput of the returned formatter escaped by \r
func (f FormatBuilder) CR() FormatBuilder {
	f.cr = true
	return f
}

// Makes the otput of the returned formatter escaped by \n
func (f FormatBuilder) LF() FormatBuilder {
	f.lf = true
	return f
}

// Makes the otput of the returned formatter escaped by \r\n
func (f FormatBuilder) CRLF() FormatBuilder {
	f.cr = true
	f.lf = true
	return f
}

// Returs a string, formatted with selected options
func (f FormatBuilder) Format(s string) string {
	if f.center {
		s = C(s)
	}
	if f.bold {
		s = B(s)
	}
	if f.cr && f.lf {
		s = CRLF(s)
	} else if f.cr {
		s = CR(s)
	} else if f.lf {
		s = LF(s)
	}
	return s
}

// Formats a string and writes it to io.Writer
func (f FormatBuilder) WriteTo(s string, w io.Writer) (int, error) {
	s = f.Format(s)
	return w.Write([]byte(s))
}

// Apllies the formatter to each individual line separately
func (f FormatBuilder) FormatLines(lines ...string) string {
	sb := strings.Builder{}

	for _, line := range lines {
		_, _ = sb.WriteString(f.Format(line))
	}

	return sb.String()
}
