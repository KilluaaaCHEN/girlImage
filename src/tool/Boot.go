package tool

import (
	"fmt"
	"time"
)

var t = time.Now()

func End() {
	elapsed := time.Since(t)
	fmt.Println("\n\nDone\nElapsed time:", elapsed)
	Wait()
}

func Wait() {
	fmt.Print("Enter any key to exit the program...")
	var enter string
	fmt.Scanln(&enter)
}

