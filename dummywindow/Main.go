package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)
	if len(os.Args) != 2 {
		log.Fatal("Expect one argument: Window title")
	} else if win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL); err != nil {
		log.Fatal("Unable to create window:", err)
	} else {
		win.SetTitle(os.Args[1])
		if cssProvider, err := gtk.CssProviderNew(); err != nil {
			log.Fatal("Unable to create cssProvider", err)
		} else {
			cssProvider.LoadFromData(`
			window {
				background-color: rgb(0,0,0);
			}
			`)
			if styleContext, err := win.GetStyleContext(); err != nil {
				log.Fatal("Unable to GetStyleContext", err)
			} else {
				styleContext.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
			}
		}


		win.Connect("destroy", func() {
			gtk.MainQuit()
		})

		win.ShowAll()
		gtk.Main()
	}
}
