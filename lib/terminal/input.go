package terminal

import (
	"context"
	"os"
	"sync"
)

// Captures a single keypress in raw mode, blocking the goroutine
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

var stdinLock sync.Mutex

// Transforms stdin into a stream of keypresses
// This function should generally be called only once and it would block on successive calls
// until the context is cancelled
//
// Note that even after the context is cancelled, the stream would process one more key
// before closing
func StdinIntoStream(ctx context.Context, n int) <-chan Key {
	stdinLock.Lock()

	c := make(chan Key, n)
	go func() {
		for ctx.Err() == nil {
			k, err := CaptupreKey()
			if err != nil {
				return
			}

			c <- k
		}

		close(c)
		stdinLock.Unlock()

	}()

	return c
}
