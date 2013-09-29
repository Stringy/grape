package grape

import (
	"net/url"
	"time"
)

type Comment struct {
	Author          string
	Body            string
	ScoreHidden     bool
	Ups             int
	Downs           int
	ApprovedBy      string `json:"approved_by"`
	AuthorFlairCSS  string `json:"author_flair_css_class"`
	AuthorFlairText string `json:"author_flair_text"`
	BannedBy        string `json:"banned_by"`
	Edited          time.Time
	Gilded          int
	Likes           bool
	LinkId          string `json:"link_id"`
	LinkTitle       string `json:"link_title"`
	NumReports      int    `json:"num_reports"`
	ParentId        string `json:"parent_id"`
	Sub             string `json:"subreddit"`
	SubId           string `json:"subreddit_id"`
	Distinguished   string
	Replies         []Comment
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
	if c.Author != user.Name {
		return incorrectOwnerError
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
