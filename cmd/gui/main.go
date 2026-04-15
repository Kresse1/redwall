package main

import (
	"fmt"
	"sync"

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

	var wg sync.WaitGroup
	var errPrev, errPosts error

	wg.Add(2)
	go func() {
		defer wg.Done()
		errPrev = controller.SavePrevious()
	}()
	go func() {
		defer wg.Done()
		errPosts = controller.LoadPosts()
	}()
	wg.Wait()
	if errPrev != nil {
		fmt.Println("Error: ", errPrev)
	}
	if errPosts != nil {
		fmt.Println("Error: ", errPosts)
	}

	w.SetContent(controller.BuildUI())
	go controller.DownloadImages()
	w.ShowAndRun()
}
