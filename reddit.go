package reddit

import (
	_ "bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reddit/client"
	"reddit/config"
	"reddit/things"
	"strings"
)

var cfg = config.GetInstance()

// GetSubreddit gets the front page of a named subreddit
// TODO: add support for arbitrary number of posts returned
func GetSubreddit(sub string) (*things.Subreddit, error) {
	//rlog.Printf("Getting subreddit: %s\n", sub)
	b, err := client.MakeGetRequest(fmt.Sprintf(cfg.Url["subreddit"], sub))
	if err != nil {
		return nil, err
	}
	rresp := new(things.RedditResponse)
	err = json.Unmarshal(b, rresp)
	rresp.Data.Name = sub
	return &rresp.Data, nil
}

// GetSubredditN gets the first n items from a subreddit
// returns a subreddit object containing the items
func GetSubredditN(sub string, n int) (*things.Subreddit, error) {
	var tempsub *things.Subreddit
	for ; n > 0; n -= 100 {
		data := url.Values{
			"limit": {fmt.Sprintf("%d", n)},
		}
		if tempsub != nil && len(tempsub.Items) != 0 {
			data.Set("after", tempsub.Items[len(tempsub.Items)-1].Name)
		}
		b, err := client.MakePostRequest(fmt.Sprintf(cfg.Url["subreddit"], sub), &data)
		if err != nil {
			return nil, err
		}
		rresp := new(things.RedditResponse)
		err = json.Unmarshal(b, rresp)
		if err != nil {
			return nil, err
		}
		if tempsub == nil {
			tempsub = &rresp.Data
		} else {
			tempsub.Items = append(tempsub.Items, rresp.Data.Items...)
		}
	}
	return tempsub, nil
}

// GetFrontPage currently gets the front page of *default* reddit
// TODO: apply this to currently logged in user
func GetFrontPage(user *things.Redditor) (*things.Subreddit, error) {
	b, err := client.MakeGetRequest(cfg.ApiUrl["frontpage"])
	if err != nil {
		return nil, err
	}
	rresp := new(things.RedditResponse)
	err = json.Unmarshal(b, rresp)
	if err != nil {
		return nil, err
	}
	return &rresp.Data, nil
}

// GetRedditor returns information about a given redditor
func GetRedditor(user string) (*things.Redditor, error) {
	// rlog.Printf("getting Redditor: %s\n", user)
	b, err := client.MakeGetRequest(fmt.Sprintf(cfg.Url["user"], user))
	if err != nil {
		return nil, err
	}
	uresp := new(things.UserResponse)
	err = json.Unmarshal(b, uresp)
	if err != nil {
		return nil, err
	}
	return &uresp.Data, nil
}

// Login logs a user into reddit through the api login page
// returns the same errors recieved from reddit, if applicable
// otherwise returns a redditor with populated modhash and cookie strings
func Login(user, pass string, rem bool) (*things.Redditor, error) {
	// rlog.Printf("logging in to user: %s\n", user)
	data := url.Values{
		"user":     {user},
		"passwd":   {pass},
		"api_type": {"json"},
		"rem":      {fmt.Sprintf("%v", rem)},
	}
	b, err := client.MakePostRequest(cfg.ApiUrl["login"], &data)
	if err != nil {
		return nil, err
	}
	loginResp := new(things.LoginResponse)
	err = json.Unmarshal(b, &loginResp)
	if err != nil {
		return nil, err
	}
	resperrs := loginResp.Json.Errors
	if len(resperrs) != 0 {
		str := ""
		for _, group := range resperrs {
			str += strings.Join(group, " ") + "\n"
		}
		return nil, errors.New("Login Error: " + str)
	}

	redditor := things.NewRedditor()
	redditor.Name = user
	redditor.ModHash = loginResp.Json.Data.ModHash
	return redditor, nil
}
