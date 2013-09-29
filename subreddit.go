package grape

import (
	"fmt"
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

func (s *Subreddit) GetUrl() string {
	return fmt.Sprintf(config.GetUrl("subreddit"), s.Name)
}

func NewSubreddit() *Subreddit {
	sub := new(Subreddit)
	sub.RWMutex = new(sync.RWMutex)
	return sub
}
