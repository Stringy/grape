package reddit

import (
	_ "bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	actual_url *url.URL
	UserAgent  = "Go Reddit API by String217 v0.1"
)

const (
	local = "reddit.local"

	home_url       = "http://reddit.local/"
	login_url      = "http://reddit.local/api/login"
	subreddit_url  = "http://reddit.local/r/%s.json"
	frontpage_url  = "http://reddit.local/.json"
	user_url       = "http://reddit.local/user/%s/about.json"
	me_url         = "http://reddit.local/api/me.json"
	comment_url    = "http://reddit.local/r/%s/%s.json"
	user_avail_url = "http://reddit.local/api/username_available.json"
	submit_url     = "http://reddit.local/api/submit"
	delete_url     = "http://reddit.local/api/delete_user"
	captcha_url    = "http://reddit.local/api/new_captcha"

	KindLink = "link"
	KindSelf = "self"
)

func init() {
	actual_url, _ = url.Parse("http://reddit.local/")
}

// GetSubreddit gets the front page of a named subreddit
// TODO: add support for arbitrary number of posts returned
func GetSubreddit(sub string) (*Subreddit, error) {
	link, _ := url.Parse(fmt.Sprintf(subreddit_url, sub))
	b, err := getJsonBytes(link)
	if err != nil {
		return nil, err
	}
	rresp := new(redditResponse)
	err = json.Unmarshal(b, rresp)
	rresp.Data.Name = sub
	return &rresp.Data, nil
}

// GetFrontPage currently gets the front page of *default* reddit
// TODO: apply this to currently logged in user
func GetFrontPage() (*Subreddit, error) {
	link, _ := url.Parse(frontpage_url)
	b, err := getJsonBytes(link)
	if err != nil {
		return nil, err
	}
	rresp := new(redditResponse)
	err = json.Unmarshal(b, rresp)
	if err != nil {
		return nil, err
	}
	return &rresp.Data, nil
}

// GetRedditor returns information about a given redditor
func GetRedditor(user string) (*Redditor, error) {
	link, _ := url.Parse(fmt.Sprintf(user_url, user))
	b, err := getJsonBytes(link)
	if err != nil {
		return nil, err
	}
	uresp := new(userResponse)
	err = json.Unmarshal(b, uresp)
	if err != nil {
		return nil, err
	}
	return &uresp.Data, nil
}

//Login logs a user into reddit through the api login page
//returns the same errors recieved from reddit, if applicable
//otherwise returns a redditor with populated modhash and cookie strings
func Login(user, pass string, rem bool) (*Redditor, error) {
	link, _ := url.Parse(login_url)
	data := url.Values{
		"user":     {user},
		"passwd":   {pass},
		"api_type": {"json"},
		"rem":      {fmt.Sprintf("%v", rem)},
	}
	b, err := getPostJsonBytes(link, &data)
	if err != nil {
		return nil, err
	}
	loginResp := new(loginResponse)
	err = json.Unmarshal(b, &loginResp)
	if err != nil {
		return nil, err
	}
	if len(loginResp.Json.Errors) != 0 {
		str := ""
		for _, group := range loginResp.Json.Errors {
			str += strings.Join(group, " ") + "\n"
		}
		return nil, errors.New("Login Error: " + str)
	}

	redditor := new(Redditor)
	redditor.Name = user
	redditor.ModHash = loginResp.Json.Data.ModHash
	return redditor, nil
}
