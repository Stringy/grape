package grape

import (
	"encoding/json"
	"fmt"
	"net/url"
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
	return Config.Host + s.Url
}

func NewSubreddit() *Subreddit {
	sub := new(Subreddit)
	sub.RWMutex = new(sync.RWMutex)
	return sub
}

// GetHot returns a list of "hot" submissions from the subreddit
func (s *Subreddit) GetHot(limit int) ([]Submission, error) {
	if limit == 0 {
		log.Printf("requested 0 hot entries from /r/%s", s.Name)
		return make([]Submission, 0), nil
	}
	data := &url.Values{
		"limit": {fmt.Sprintf("%d", limit)},
	}
	b, err := makePostRequest(fmt.Sprintf(Config.GetUrl("hot"), s.Name), data)
	if err != nil {
		return nil, err
	}
	ts := new(redditResponse)
	err = json.Unmarshal(b, ts)
	if err != nil {
		return nil, err
	}
	items := make([]Submission, len(ts.Data.Items))
	for i, subm := range ts.Data.Items {
		items[i] = subm.Submission
	}
	return items, nil
}

// GetTop returns a list of "top" submissions for a particular time
// period (hour, day, month, year, all)
func (s *Subreddit) GetTop(t period, limit int) ([]Submission, error) {
	if limit == 0 {
		log.Printf("requested 0 top entries from /r/%s in period %s", s.Name, t)
		return make([]Submission, 0), nil
	}
	data := &url.Values{
		"limit": {fmt.Sprintf("%d", limit)},
		"t":     {fmt.Sprintf("%d", t)},
	}
	b, err := makePostRequest(fmt.Sprintf(Config.GetUrl("top"), s.Name), data)
	if err != nil {
		return nil, err
	}
	ts := new(redditResponse)
	err = json.Unmarshal(b, ts)
	if err != nil {
		return nil, err
	}
	items := make([]Submission, len(ts.Data.Items))
	for i, subm := range ts.Data.Items {
		items[i] = subm.Submission
	}
	return items, nil
}

// GetNew returns a list of "new" submissions from the subreddit
func (s *Subreddit) GetNew(limit int) ([]Submission, error) {
	if limit == 0 {
		log.Printf("requested 0 new entries from /r/%s", s.Name)
		return make([]Submission, 0), nil
	}
	data := &url.Values{
		"limit": {fmt.Sprintf("%d", limit)},
	}
	b, err := makePostRequest(fmt.Sprintf(Config.GetUrl("new"), s.Name), data)
	if err != nil {
		return nil, err
	}
	ts := new(redditResponse)
	err = json.Unmarshal(b, ts)
	if err != nil {
		return nil, err
	}
	items := make([]Submission, len(ts.Data.Items))
	for i, subm := range ts.Data.Items {
		items[i] = subm.Submission
	}
	return items, nil
}
