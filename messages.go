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

// MarkAsRead marks the messages as read for the user
func (m *Message) MarkAsRead(user *Redditor) error {
	return nil
}

// MarkAsUnread marks the message as unread for the user
func (m *Message) MarkAsUnread(user *Redditor) error {
	return nil
}

// Block blocks the author of the message for the user
func (m *Message) Block(user *Redditor) error {
	return nil
}
