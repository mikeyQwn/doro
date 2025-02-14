package main

import (
	"fmt"
	"github.com/mikeyQwn/doro/lib"
	"time"
)

func main() {
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
