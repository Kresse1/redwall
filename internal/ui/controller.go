package ui

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Kresse1/redwall/internal/imaging"
	"github.com/Kresse1/redwall/internal/reddit"
	"github.com/Kresse1/redwall/internal/wallpaper"
)

type Controller struct {
	client      *reddit.Client
	kde         *wallpaper.KDESetter
	screen      *imaging.Screen
	posts       []reddit.Post
	images      [][]byte
	selectedID  int
	previousWP  string
	preview     *canvas.Image
	currentSlot bool
	cacheDir    string
}

func NewController(
	client *reddit.Client,
	kde *wallpaper.KDESetter,
	screen *imaging.Screen,
) (*Controller, error) {

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(cacheDir, "redwall")
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	return &Controller{
		client:   client,
		kde:      kde,
		screen:   screen,
		cacheDir: path,
	}, nil

}

func (c *Controller) SavePrevious() error {
	prev, err := c.kde.Current()
	if err != nil {
		return err
	}
	c.previousWP = prev
	return nil
}

func (c *Controller) LoadPosts() error {
	posts, err := c.client.FetchPosts()
	if err != nil {
		return err
	}
	c.posts = posts
	return nil
}

func (c *Controller) DownloadImages() error {
	var wg sync.WaitGroup
	errs := make([]error, len(c.posts))
	c.images = make([][]byte, len(c.posts))

	for i, post := range c.posts {
		wg.Add(1)
		go func(i int, post reddit.Post) {
			defer wg.Done()
			img, err := c.client.DownloadImage(post.URL)
			if err != nil {
				errs[i] = err
				return
			}
			c.images[i] = img
		}(i, post)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) SelectPost(id int) {
	if c.images[id] == nil {
		return
	}
	c.selectedID = id
	img, _, err := image.Decode(bytes.NewReader(c.images[id]))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	c.preview.Image = img
	c.preview.Refresh()
}

func (c *Controller) SetWallpaper() error {

	c.currentSlot = !c.currentSlot
	scaledImageBytes, err := imaging.Scale(c.screen.Width, c.screen.Height, bytes.NewReader(c.images[c.selectedID]))
	if err != nil {
		return err
	}
	var path string
	if c.currentSlot {
		path = filepath.Join(c.cacheDir, "slot0.jpg")
	} else {
		path = filepath.Join(c.cacheDir, "slot1.jpg")
	}

	if err := os.WriteFile(path, scaledImageBytes, 0644); err != nil {
		return err
	}

	return c.kde.Set(path)
}
func (c *Controller) Undo() error {
	return c.kde.Set(c.previousWP)
}

func (c *Controller) BuildUI() fyne.CanvasObject {
	leftPanel := widget.NewList(
		func() int { return len(c.posts) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			label := item.(*widget.Label)
			label.SetText(c.posts[id].Title)
		},
	)
	leftPanel.OnSelected = c.SelectPost

	buttons := container.NewHBox(
		widget.NewButton("Set Wallpaper", func() {
			if err := c.SetWallpaper(); err != nil {
				fmt.Println("Error:", err)
			}
		}),
		widget.NewButton("Undo", func() {
			if err := c.Undo(); err != nil {
				fmt.Println("Error:", err)
			}
		}),
	)

	c.preview = canvas.NewImageFromImage(nil)
	c.preview.FillMode = canvas.ImageFillContain

	rightPanel := container.NewBorder(nil, buttons, nil, nil, c.preview)
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.3

	return split

}
