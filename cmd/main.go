package main

import (
	"github.com/mikeyQwn/doro/bin"
	"github.com/mikeyQwn/doro/lib/terminal"
	"log"
)

func main() {
	restore, err := terminal.IntoRaw()
	if err != nil {
		log.Fatal(err)
	}
	defer restore()

	bin.Run()
}
