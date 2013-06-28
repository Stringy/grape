package reddit

type Comment struct {
	Author      string
	Body        string
	ScoreHidden bool
	Ups         int
	Downs       int
	Replies     []Comment
}

func (c *Comment) Reply(user *Redditor, body string) (*Comment, error) {
	return nil, nil
}

func (c *Comment) Edit(user *Redditor, body string) (*Comment, error) {
	return nil, nil
}

func (c *Comment) Report() error {
	return nil
}

func (c *Comment) Hide() error {
	return nil
}

func (c *Comment) Unhide() error {
	return nil
}

func (c *Comment) MarkNsfw() error {
	return nil
}

func (c *Comment) UnmarkNsfw() error {
	return nil
}

func (c *Comment) Save() (*Comment, error) {
	return nil, nil
}

func (c *Comment) MoreChildren() ([]*Comment, error) {
	return nil, nil
}

// true for up, false for down
func (c *Comment) Vote(vote bool) error {
	return nil
}
