package reddit

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/bitly/go-simplejson"
	"io"
	_ "io/ioutil"
	"net/http"
	_ "time"
)

type Redditor struct {
	Name     string
	LKarma   int  `json:"link_karma"`
	CKarma   int  `json:"comment_karma"`
	IsFriend bool `json:"is_friend"`
	HasMail  bool `json:"has_mail"`
	IsOver18 bool `json:"over_18"`
	IsGold   bool `json:"is_gold"`
	IsMod    bool `json:"is_mod"`
}

type UserResponse struct {
	Data Redditor
}

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
	cresp := make([]*CommentsResponse, 2)
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

type Subreddit struct {
	Id    string
	Name  string
	Items []struct {
		RedditPost `json:"data"`
	} `json:"children"`
}

type RedditResponse struct {
	Data Subreddit
}

type jsonComment struct {
	Author      string
	Body        string
	ScoreHidden bool `json:"score_hidden"`
	Ups         int
	Downs       int
	Replies     struct {
		Data struct {
			Children []struct {
				Data jsonComment
			}
		}
	}
}

type Comment struct {
	Author      string
	Body        string
	ScoreHidden bool
	Ups         int
	Downs       int
	Replies     []Comment
}

func commentFromJson(jComm jsonComment) Comment {
	comment := new(Comment)
	comment.Author = jComm.Author
	comment.Body = jComm.Body
	comment.ScoreHidden = jComm.ScoreHidden
	comment.Ups = jComm.Ups
	comment.Downs = jComm.Downs
	comment.Replies = make([]Comment, len(jComm.Replies.Data.Children))
	for i, jCommReply := range jComm.Replies.Data.Children {
		comment.Replies[i] = commentFromJson(jCommReply.Data)
	}
	return *comment
}

type CommentsResponse struct {
	Data struct {
		Children []struct {
			Data jsonComment `json:"data"`
		}
	}
}
