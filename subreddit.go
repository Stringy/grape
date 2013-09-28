package reddit

import (
	"sync"
)

type Subreddit struct {
	Id    string
	Name  string
	Items []struct {
		Submission `json:"data"`
	} `json:"children"`
	*sync.RWMutex
}

func NewSubreddit() *Subreddit {
	sub := new(Subreddit)
	sub.RWMutex = new(sync.RWMutex)
	return sub
}
