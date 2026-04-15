package main

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/Kresse1/redwall/internal/imaging"
	"github.com/Kresse1/redwall/internal/reddit"
	"github.com/Kresse1/redwall/internal/wallpaper"
)

func main() {
	client := reddit.NewClient()
	kdeSetter := wallpaper.NewKDESetter()
	screen, err := imaging.NewScreen()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Test Screen
	fmt.Printf("Screen size %dx%d\n", screen.Width, screen.Height)

	posts, err := client.FetchPosts()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i, post := range posts {
		fmt.Printf("[%d] %s - %s\n", i, post.Title, post.URL)
	}

	fmt.Println("Choose a wallpaper")
	var postNumber int
	_, err = fmt.Scanln(&postNumber)
	if err != nil {
		fmt.Println("Eroor:", err)
		return
	}
	if postNumber < 0 || postNumber >= len(posts) {
		fmt.Printf("Post %d not found\n", postNumber)
		return
	}

	downloadUrl := posts[postNumber].URL
	parsedUrl, err := url.Parse(downloadUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dataFormat := path.Ext(parsedUrl.Path)
	imgBytes, err := client.DownloadImage(downloadUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var imagePath = "image" + dataFormat
	scaledImgBytes, err := imaging.Scale(screen.Width, screen.Height, bytes.NewReader(imgBytes))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = os.WriteFile(imagePath, scaledImgBytes, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	currentWallpaper, err := kdeSetter.Current()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	absFilePath, err := filepath.Abs(imagePath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = kdeSetter.Set(absFilePath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	println("Keep new Wallpaper? (yes)")
	var confirm string
	_, err = fmt.Scanln(&confirm)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	if confirm == "yes" {
		return
	} else {
		err = kdeSetter.Set(currentWallpaper)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

	}
}
