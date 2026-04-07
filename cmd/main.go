package main

import (
	"fmt"
	"github.com/Kresse1/redwall/internal/reddit"
)

func main() {
	client := reddit.NewClient()
	posts, err := client.FetchPosts()
	if err != nil {
		fmt.Println("Error:", err)
	}
	for _, post := range posts {
		fmt.Println(post.Title, "-", post.URL)
	}
}
