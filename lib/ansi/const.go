// ANSI escape codes

package ansi

// A list of ANSI escape codes ready to use in strings
const (
	RESET         = "\033[0m"
	BOLD          = "\033[1m"
	ERASE_LINE    = "\033[2K"
	UP_LINE       = "\033[A"
	UP_START_LINE = "\033[F"

	RED = "\033[31m"

	CARRIAGE_RETURN = "\r"
	LINE_FEED       = "\n"
)
