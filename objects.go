package reddit

import (
	"encoding/json"
	"fmt"
	_ "github.com/bitly/go-simplejson"
	_ "io/ioutil"
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
	cresp := make([]CommentsResponse, 2)
	err = json.NewDecoder(resp.Body).Decode(cresp)
	fmt.Println(cresp)
	comments := make([]Comment, len(cresp[1].Data.Children))
	for i, comment := range cresp[1].Data.Children {
		comments[i] = comment.Data
	}
	fmt.Println(cresp)
	//bytes, _ := ioutil.ReadAll(resp.Body)
	//_, _ = sjson.NewJson(bytes)
	//fmt.Printf("%s", bytes)
	return comments
	//return nil
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
	Replies     []struct {
		Data struct {
			Children []struct {
				Data Comment
			}
		}
	}
}

type CommentsResponse struct {
	Data struct {
		Children []struct {
			Data Comment `json:"data"`
		}
	}
}
