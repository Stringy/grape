package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Submission struct {
	Title       string
	Url         string
	NumComments int `json:"num_comments"`
	Author      string
	IsSelf      bool `json:"is_self"`
	IsNSFW      bool `json:"over_18"`
	SelfText    string
	Created     float64 `json:"created_utc"`
	Score       int
	Ups         int
	Downs       int
	Sub         string `json:"subreddit"`
	*Thing
}

func (r *Submission) String() string {
	str := fmt.Sprintf(
		"Title: %s\n\t%d Up \n\t%d Down\n\tAuthor: %s\n\tSub: %s\n",
		r.Title,
		r.Ups,
		r.Downs,
		r.Author,
		r.Sub)
	return str
}

func (r *Submission) GetComments() []Comment {
	b, err := makeGetRequest(fmt.Sprintf(Urls["comment"], r.Sub, r.Id))
	if err != nil {
		panic(err)
	}
	cresp := make([]*commentsResponse, 2)
	err = json.Unmarshal(b, &cresp)
	comments := make([]Comment, len(cresp[1].Data.Children))
	for i, comment := range cresp[1].Data.Children {
		comments[i] = commentFromJson(comment.Data)
	}
	return comments
}

func (r *Submission) PostComment(user *Redditor, body string) error {
	if !user.IsLoggedIn() {
		return NotLoggedInError
	}
	data := &url.Values{
		"api_type": {"json"},
		"text":     {body},
		"uh":       {user.ModHash},
		"thing_id": {"t6_" + r.Id},
	}
	b, err := makePostRequest(ApiUrls["comment"], data)
	if err != nil {
		return err
	}
	errstruct := new(struct {
		Json struct {
			Errors [][]string
		}
	})
	err = json.Unmarshal(b, &errstruct)
	if err != nil {
		return err
	}
	if len(errstruct.Json.Errors) != 0 {
		return errors.New(strings.Join(errstruct.Json.Errors[0], ", "))
	}
	return nil
}
