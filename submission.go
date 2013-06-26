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
	link, _ := url.Parse(fmt.Sprintf(comment_url, r.Sub, r.Id))
	b, err := getJsonBytes(link)
	if err != nil {
		panic(err)
	}
	cresp := make([]*commentsResponse, 2)
	err = json.Unmarshal(b, &cresp)
	//	fmt.Println(cresp[1])
	comments := make([]Comment, len(cresp[1].Data.Children))
	for i, comment := range cresp[1].Data.Children {
		comments[i] = commentFromJson(comment.Data)
	}
	return comments
}
