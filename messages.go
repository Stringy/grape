package grape

type Message struct {
	WasComment bool
	Body       string
	Subject    string
	Subreddit  string // if comment
	ParentId   string `json:"parent_id"`
	New        bool
	Author     string
	Recipient  string `json:"dest"`
	Context    string
	Likes      bool
	LinkTitle  string `json:"link_title"` // if comment
	Replies    string
	*Thing
}
