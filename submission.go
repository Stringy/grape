package reddit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	req := constructDefaultRequest(
		"GET",
		fmt.Sprintf(comment_url, r.Sub, r.Id))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(err)
	}
	cresp := make([]*commentsResponse, 2)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf.Bytes(), &cresp)
	//	fmt.Println(cresp[1])
	comments := make([]Comment, len(cresp[1].Data.Children))
	for i, comment := range cresp[1].Data.Children {
		comments[i] = commentFromJson(comment.Data)
	}
	return comments
}
