package grape

type Message struct {
	WasComment bool
	Body       string
	Subject    string
	Subreddit  string
	ParentId   string `json:"parent_id"`
	New        bool
	Author     string
	Recipient  string `json:"dest"`
	*Thing
}
