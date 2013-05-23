package reddit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	client *http.Client //default http client for requests
)

const (
	home_url       = "http://www.reddit.com/"
	login_url      = "http://www.reddit.com/api/login"
	subreddit_url  = "http://www.reddit.com/r/%s.json"
	frontpage_url  = "http://www.reddit.com/.json"
	user_url       = "http://www.reddit.com/user/%s/about.json"
	me_url         = "http://www.reddit.com/api/me.json"
	comment_url    = "http://www.reddit.com/r/%s/%s.json"
	user_avail_url = "http://www.reddit.com/api/username_available.json"

	UserAgent = "Go Reddit API by String217 v0.1"
)

func init() {
	client = new(http.Client)
}

// GetSubreddit gets the front page of a named subreddit
// TODO: add support for arbitrary number of posts returned
func GetSubreddit(sub string) *Subreddit {
	req := constructDefaultRequest(
		"GET",
		fmt.Sprintf(subreddit_url, sub))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rresp := new(redditResponse)
	err = json.NewDecoder(resp.Body).Decode(rresp)
	rresp.Data.Name = sub
	return &rresp.Data
}

// GetFrontPage currently gets the front page of *default* reddit
// TODO: apply this to currently logged in user
func GetFrontPage() *Subreddit {
	req := constructDefaultRequest(
		"GET",
		frontpage_url)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rresp := new(redditResponse)
	err = json.NewDecoder(resp.Body).Decode(rresp)
	return &rresp.Data
}

// GetRedditor returns information about a given redditor
func GetRedditor(user string) (*Redditor, error) {
	req := constructDefaultRequest(
		"GET",
		fmt.Sprintf(user_url, user))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	uresp := new(userResponse)
	err = json.NewDecoder(resp.Body).Decode(uresp)
	if err != nil {
		return nil, err
	}
	return &uresp.Data, nil
}

//Login logs a user into reddit through the api login page
//returns the same errors recieved from reddit, if applicable
//otherwise returns a redditor with populated modhash and cookie strings
func Login(user, pass string, rem bool) (*Redditor, error) {
	resp, err := http.PostForm(
		login_url,
		url.Values{
			"user":     {user},
			"passwd":   {pass},
			"api_type": {"json"},
			"rem":      {fmt.Sprintf("%v", rem)},
		})
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		panic(err)
	}
	loginResp := new(loginResponse)
	err = json.Unmarshal(buf.Bytes(), &loginResp)
	if err != nil {
		panic(err)
	}
	fmt.Println(loginResp)
	//	fmt.Println(buf.String())
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
	redditor.Cookie = loginResp.Json.Data.Cookie
	return redditor, nil
}

func usernameAvailable(user string) bool {
	resp, err := http.Get(user_avail_url)
	//url.Values{
	// 	"user": {user},
	// }
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	return false
}

func constructDefaultRequest(request_type, url string) *http.Request {
	req, err := http.NewRequest(request_type, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-agent", UserAgent)
	return req
}
