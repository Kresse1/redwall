package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Kresse1/redwall/internal/imaging"
	"github.com/Kresse1/redwall/internal/reddit"
	"github.com/Kresse1/redwall/internal/ui"
	"github.com/Kresse1/redwall/internal/wallpaper"
)

func main() {

	screen, err := imaging.NewScreen()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	controller, err := ui.NewController(reddit.NewClient(), wallpaper.NewKDESetter(), screen)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	a := app.New()
	w := a.NewWindow("redwall")
	w.Resize(fyne.NewSize(800, 600))
	if err := controller.SavePrevious(); err != nil {
		fmt.Println("Error", err)
		return
	}
	if err := controller.LoadPosts(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := controller.DownloadImages(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	w.SetContent(controller.BuildUI())
	w.ShowAndRun()
}
