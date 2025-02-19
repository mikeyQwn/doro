// Time formatting functions

package lib

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
