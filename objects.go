package reddit

import (
	"fmt"
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
	Comments    []Comment
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

type Comment struct {
	Author      string
	Body        string
	ScoreHidden bool `json:"score_hidden"`
	Ups         int
	Downs       int
	Replies     []Comment
}

type CommentsResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Children []struct {
					Data Comment
				}
			}
		}
	}
}
