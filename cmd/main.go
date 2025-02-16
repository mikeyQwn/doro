package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mikeyQwn/doro/lib"
	"github.com/mikeyQwn/doro/lib/terminal"
)

func main() {
	res, err := terminal.IntoRaw()
	if err != nil {
		log.Fatal(err)
	}
	k, err := terminal.CaptupreKey()
	if k == terminal.KEY_ARROW_LEFT {
		fmt.Println("LEFT")
	}
	res()
	log.Fatal()

	fmt.Print("Doing stuff... ")
	spinnerFrames := []string{
		"\\ Here",
		"| We",
		"/ Go",
		"- ...",
		"\\ Almost",
		"| Done",
		"/ ...",
		"- ...",
	}

	a := lib.NewAnimation()

	for {
		for _, c := range spinnerFrames {
			fmt.Print(a.NextFrame(c))
			time.Sleep(time.Second)
		}
	}
}
