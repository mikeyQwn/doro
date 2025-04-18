// Time formatting functions

package bin

import (
	"strconv"
	"time"
)

// Formats duration, rounding down to nearest minute, if necessary
func FormatDurationMinPrec(d time.Duration) string {
	if h := int(d.Hours()); h > 0 && d.Hours() == float64(h) {
		return strconv.Itoa(h) + "h"
	}

	if m := int(d.Minutes()); m > 0 {
		return strconv.Itoa(m) + "m"
	}

	return strconv.Itoa(int(d.Seconds())) + "s"
}

// Converts float representation of percentage
// To a string with a single decimal precision, e. g. 0.4512 -> 45.1%
func FormatPercent(p float64) string {
	return strconv.FormatFloat(p*100, 'f', 1, 64) + "%"
}

func FormatTimer(elapsed, total time.Duration) string {
	return strconv.Itoa(int(elapsed.Minutes())) + "/" + strconv.Itoa(int(total.Minutes())) + "m"
}
