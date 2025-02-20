package input

import (
	"context"
	"os"
	"sync"
)

type KeyStream <-chan Key

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
func StdinIntoStream(ctx context.Context, n int) KeyStream {
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

// Waits until user presses the given `key` or the stream closes
func (s KeyStream) WaitKey(key Key) {
	for k := range s {
		if k == key {
			return
		}
	}
}

// Transforms the `KeyStream` so that it applies a function to
// every key received.
//
// Note that this call gives up "ownership" of the `KeyStream`
// and the returned `KeyStream` should be used instead
func (s KeyStream) Map(f func(Key) Key) KeyStream {
	newS := make(chan Key, cap(s))
	go func() {
		for k := range s {
			newS <- f(k)
		}
		close(newS)
	}()

	return newS
}

// Transforms the `KeyStream` so that it triggers a handler
// when `CtrlC` is pressed
//
// Note that this call gives up "ownership" of the `KeyStream`
// and the returned `KeyStream` should be used instead
func (s KeyStream) HandleCtrlC(handler func()) KeyStream {
	return s.
		Map(func(k Key) Key {
			if k == KEY_CTRL_C {
				handler()
			}

			return k
		})
}
