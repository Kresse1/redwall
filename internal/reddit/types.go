package reddit

type Response struct {
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	After    string  `json:"after"`
	Children []Child `json:"children"`
}

type Child struct {
	Data Post `json:"data"`
}

type Post struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	URL      string  `json:"url"`
	PostHint string  `json:"post_hint"`
	Score    int     `json:"score"`
	Over18   bool    `json:"over_18"`
	Preview  Preview `json:"preview"`
}

type Preview struct {
	Images []PreviewImage `json:"images"`
}

type PreviewImage struct {
	Source      ImageSource   `json:"source"`
	Resolutions []ImageSource `json:"resolutions"`
}

type ImageSource struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
