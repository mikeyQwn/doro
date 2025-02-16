package terminal

import (
	"context"
	"fmt"
	"os"
)

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

	fmt.Printf("RUNE IS: \"%c\"", rune(k))

	fmt.Println(buf[:n], k, KEY_ARROW_LEFT)

	return k, nil
}

func CaptureInputStream(ctx context.Context, n int) <-chan Key {
	c := make(chan Key, n)

	// TODO: implement me

	return c
}
