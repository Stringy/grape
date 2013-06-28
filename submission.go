package reddit

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type RedditPost struct {
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
	Id          string
	Sub         string `json:"subreddit"`
}

func (r *RedditPost) String() string {
	str := fmt.Sprintf(
		"Title: %s\n\t%d Up \n\t%d Down\n\tAuthor: %s\n\tSub: %s\n",
		r.Title,
		r.Ups,
		r.Downs,
		r.Author,
		r.Sub)
	return str
}

func (r *RedditPost) GetComments() []Comment {
	b, err := getJsonBytes(fmt.Sprintf(Urls["comment"], r.Sub, r.Id))
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

func (r *RedditPost) PostComment(user *Redditor, comment *Comment) error {
	if !user.IsLoggedIn() {
		return NotLoggedInError
	}
	data := &url.Values{
		"api_type": {"json"},
		"text":     {comment.Body},
		"uh":       {user.ModHash},
		"thing_id": {"t6_" + r.Id},
	}
	_, err := getPostJsonBytes(ApiUrls["comment"], data)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedditPost) Report() error {
	return nil
}

func (r *RedditPost) MarkNsfw() error {
	return nil
}

func (r *RedditPost) UnmarkNsfw() error {
	return nil
}

func (r *RedditPost) Hide() error {
	return nil
}

func (r *RedditPost) Unhide() error {
	return nil
}

func (r *RedditPost) Save() error {
	return nil
}

func (r *RedditPost) Vote(vote bool) error {
	return nil
}

func (r *RedditPost) Info() error {
	return nil
}
