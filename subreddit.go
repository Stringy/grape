package reddit

type Subreddit struct {
	Id    string
	Name  string
	Items []struct {
		RedditPost `json:"data"`
	} `json:"children"`
}
