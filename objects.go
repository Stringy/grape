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
	Title    string
	Url      string
	Comments int `json:"num_comments"`
	Author   string
	IsSelf   bool `json:"is_self"`
	IsNSFW   bool `json:"over_18"`
	SelfText string
	Created  float64 `json:"created_utc"`
	Score    int
	Ups      int
	Downs    int
}

func (r *RedditPost) String() string {
	str := ""
	str = fmt.Sprintf("Title: %s\n", r.Title)
	str = fmt.Sprintf("%sUrl: %s\n", str, r.Url)
	str = fmt.Sprintf("%sAuthor: %s\n", str, r.Author)
	str = fmt.Sprintf("%sSelf? %v\n", str, r.IsSelf)
	str = fmt.Sprintf("%sComments: %d\n", str, r.Comments)
	str = fmt.Sprintf("%sNSFW? %v\n", str, r.IsNSFW)
	str = fmt.Sprintf("%sText: %s\n", str, r.SelfText)
	return str
}

type Subreddit struct {
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
