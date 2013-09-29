package grape

import (
	"sync"
)

type Subreddit struct {
	AccountsActive int `json:"accounts_active"`
	HideDuration   int `json:"comment_score_hide_mins"`
	Description    string
	DisplayName    string `json:"display_name"`
	HeaderImage    string `json:"header_img"`
	HeaderSize     []int  `json:"header_size"`
	HeaderTitle    string `json:"header_title"`
	Over18         bool   `json:"over_18"`
	PublicDesc     string `json:"public_description"`
	PublicTraffic  bool   `json:"public_traffic"`
	Subscribers    int
	SubmissionType string `json:"submission_type"`
	SubredditType  string `json:"subreddit_type"`
	Title          string
	Url            string // relative e.g. /r/pics
	UserIsBanned   bool   `json:"user_is_banned"`
	UserIsContrib  bool   `json:"user_is_contributor"`
	UserIsMod      bool   `json:"user_is_moderator"`
	UserIsSub      bool   `json:"user_is_subscriber"`
	Items          []struct {
		Submission `json:"data"`
	} `json:"children"`
	*sync.RWMutex
	*Thing
}

func (s *Subreddit) GetUrl() string {
	return config.Host + s.Url
}

func NewSubreddit() *Subreddit {
	sub := new(Subreddit)
	sub.RWMutex = new(sync.RWMutex)
	return sub
}
