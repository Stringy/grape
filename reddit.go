package reddit

import (
	_ "bufio"
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
	client     *http.Client //default http client for requests
	actual_url *url.URL
)

const (
	UserAgent = "Go Reddit API by String217 v0.1"

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

	KindLink = "link"
	KindSelf = "self"
)

func init() {
	client = new(http.Client)
	actual_url, _ = url.Parse("http://reddit.local/")
}

// GetSubreddit gets the front page of a named subreddit
// TODO: add support for arbitrary number of posts returned
func GetSubreddit(sub string) (*Subreddit, error) {
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
	return &rresp.Data, nil
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
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(user_url, user), nil)
	if err != nil {
		return nil, err
	}
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
	//	fmt.Println(buf.String())
	loginResp := new(loginResponse)
	err = json.Unmarshal(buf.Bytes(), &loginResp)
	if err != nil {
		panic(err)
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
	redditor.Cookies = resp.Cookies()
	return redditor, nil
}

func usernameAvailable(user string) bool {
	resp, err := http.Get(user_avail_url)
	// url.Values{
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
