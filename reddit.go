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
	client *http.Client
)

const (
	login       = "http://www.reddit.com/api/login"
	subreddit   = "http://www.reddit.com/r/%s.json"
	frontpage   = "http://www.reddit.com/.json"
	user_url    = "http://www.reddit.com/user/%s/about.json"
	me_url      = "http://www.reddi.com/api/me.json"
	comment_url = "http://www.reddit.com/r/%s/%s.json"

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
	rresp := new(redditResponse)
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
	rresp := new(redditResponse)
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
	uresp := new(userResponse)
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

func Login(user, pass string, rem bool) (*Redditor, error) {
	resp, err := http.PostForm(
		login,
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
