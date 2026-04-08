package main

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/Kresse1/redwall/internal/reddit"
)

func main() {
	client := reddit.NewClient()
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
	if postNumber < 0 || postNumber > len(posts) {
		fmt.Printf("Post %d not found\n", postNumber)
		return
	}

	downloadUrl := posts[postNumber].URL
	parsedUrl, err := url.Parse(downloadUrl)
	dataFormat := path.Ext(parsedUrl.Path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	image, err := client.DownloadImage(downloadUrl)
	if err != nil {
		fmt.Println("Error:", err)
	}
	err = os.WriteFile("image"+dataFormat, image, 0644)
	if err != nil {
		fmt.Println("Error:", err)
	}

}
