package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/Kresse1/redwall/internal/reddit"
	"github.com/Kresse1/redwall/internal/wallpaper"
)

func main() {
	client := reddit.NewClient()
	kdeSetter := wallpaper.NewKDESetter()
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
	image, err := client.DownloadImage(downloadUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var imagePath = "image" + dataFormat
	err = os.WriteFile(imagePath, image, 0644)
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
