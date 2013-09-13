package things

import (
	"net/url"
	"reddit/client"
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
	_, err := client.MakePostRequest(cfg.ApiUrl["comment"], data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) Edit(user *Redditor, body string) error {
	if !user.IsLoggedIn() {
		return NotLoggedInError
	}
	data := &url.Values{
		"api_type": {"json"},
		"text":     {body},
		"uh":       {user.ModHash},
		"thing_id": {"t6_" + c.Id},
	}
	_, err := client.MakePostRequest(cfg.ApiUrl["editusertext"], data)
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
