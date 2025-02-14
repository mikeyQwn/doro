package lib

import (
	"fmt"
	"strings"
	"unicode"
)

// The state of a alternating line of text
//
// To use this, `NextFrame` should be called for every frame of the animation
// and produced string is supposed to be written to something that supports escape codes
type Animation struct {
	prevLen int
}

func NewAnimation() Animation {
	return Animation{
		prevLen: 0,
	}
}

// Generates a string that replaces text returned by the previous call using esape codes.
// Make sure to not use a function that adds newline after passed string to ouput the `NextFrame`,
// so don't use fmt.Println(a.NextFrame("your input here")), fmt.Print() should be preferred
//
// # Panics
//
// Panics when the given `line` contains escape codes
func (a *Animation) NextFrame(line string) string {
	if hasEscapes(line) {
		panic("line provided to `NextFrame` contains escape codes")
	}

	backs := strings.Repeat("\b", a.prevLen)
	spaces := strings.Repeat(" ", a.prevLen)

	res := fmt.Sprintf("%s%s%s%s", backs, spaces, backs, line)
	a.prevLen = len(line)

	return res
}

func hasEscapes(line string) bool {
	return strings.ContainsFunc(line, unicode.IsControl)
}
