package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	subreddit  string
	userAgent  string
}

func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
		subreddit:  "wallpaper",
		userAgent:  "reddwall/1.0 (github.com/Kresse1/redwall)",
	}
}

func (c *Client) FetchPosts() ([]Post, error) {
	requestURL := fmt.Sprintf("https://www.reddit.com/r/%s/hot.json", c.subreddit)
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", c.userAgent)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var redditResponse Response
	err = json.Unmarshal(body, &redditResponse)
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, child := range redditResponse.Data.Children {
		posts = append(posts, child.Data)
	}
	return posts, nil

}
