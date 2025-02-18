package terminal

import (
	"os"
)

// Captures a single keypress in raw mode
// Returns an associated key or error if read from stdin fails
func CaptupreKey() (Key, error) {
	buf := [4]byte{0}

	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		return KEY_UNKNOWN, err
	}

	k := Key(0)
	for shift, b := range buf[:n] {
		mask := Key(b) << (shift * 8)
		k += mask
	}

	return k, nil
}
