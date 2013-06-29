package reddit

import (
	"net/url"
)

type Comment struct {
	Author string
	Body   string
	//	Id          string
	ScoreHidden bool
	Ups         int
	Downs       int
	Replies     []Comment
	*Thing
}

func (c *Comment) Reply(user *Redditor, body string) error {
	if !user.IsLoggedIn() {
		return NotLoggedInError
	}
	data := &url.Values{
		"api_type": {"json"},
		"text":     {body},
		"uh":       {user.ModHash},
		"thing_id": {"t6_" + c.Id},
	}
	_, err := getPostJsonBytes(ApiUrls["comment"], data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) Edit(user *Redditor, body string) error {
	return nil
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
