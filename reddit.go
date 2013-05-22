package reddit

import (
	"encoding/json"
	_ "errors"
	"fmt"
	_ "github.com/bitly/go-simplejson"
	_ "io/ioutil"
	"net/http"
)

var (
	client *http.Client
)

const (
	login     = "http://www.reddit.com/api/login"
	subreddit = "http://www.reddit.com/r/%s.json"
	frontpage = "http://www.reddit.com/.json"
	user_url  = "http://www.reddit.com/user/%s/about.json"

	default_ua = "Go Reddit API by String217 v0.1"
)

func init() {
	client = new(http.Client)
}

func GetSubreddit(sub string) *Subreddit {
	req := constructDefaultRequest(
		"GET",
		fmt.Sprintf(subreddit, sub))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rresp := new(RedditResponse)
	err = json.NewDecoder(resp.Body).Decode(rresp)
	rresp.Data.Name = sub
	return &rresp.Data
}

func GetFrontPage() *Subreddit {
	req := constructDefaultRequest(
		"GET",
		frontpage)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rresp := new(RedditResponse)
	err = json.NewDecoder(resp.Body).Decode(rresp)
	return &rresp.Data
}

func GetRedditor(user string) (*Redditor, error) {
	req := constructDefaultRequest(
		"GET",
		fmt.Sprintf(user_url, user))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	uresp := new(UserResponse)
	err = json.NewDecoder(resp.Body).Decode(uresp)
	if err != nil {
		return nil, err
	}
	return &uresp.Data, nil
}

func constructDefaultRequest(request_type, url string) *http.Request {
	req, err := http.NewRequest(request_type, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", default_ua)
	return req
}
