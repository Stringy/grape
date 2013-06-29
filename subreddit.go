package reddit

type Subreddit struct {
	Id    string
	Name  string
	Items []struct {
		Submission `json:"data"`
	} `json:"children"`
}
