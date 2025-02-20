package main

import (
	"github.com/mikeyQwn/doro/bin"
)

func main() {
	if err := bin.Run(); err != nil {
		bin.Fatal(err.Error())
	}
}
