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
		return notLoggedInError
	}
	data := &url.Values{
		"api_type": {"json"},
		"text":     {body},
		"uh":       {user.ModHash},
		"thing_id": {"t6_" + c.Id},
	}
	_, err := makePostRequest(config.GetApiUrl("comment"), data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) Edit(user *Redditor, body string) error {
	if !user.IsLoggedIn() {
		return notLoggedInError
	}
	data := &url.Values{
		"api_type": {"json"},
		"text":     {body},
		"uh":       {user.ModHash},
		"thing_id": {"t6_" + c.Id},
	}
	_, err := makePostRequest(config.GetApiUrl("editusertext"), data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) MoreChildren() ([]*Comment, error) {
	return nil, nil
}

// true for up, false for down
func (c *Comment) Vote(vote bool) error {
	return nil
}
