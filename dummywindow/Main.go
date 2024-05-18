package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Expect one argument: Window title")
	} else {
		a := app.New()
		w := a.NewWindow(os.Args[1])

		w.ShowAndRun()
	}
}
