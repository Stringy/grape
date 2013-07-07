package reddit

import (
	_ "bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// GetSubreddit gets the front page of a named subreddit
// TODO: add support for arbitrary number of posts returned
func GetSubreddit(sub string) (*Subreddit, error) {
	b, err := getJsonBytes(fmt.Sprintf(Urls["subreddit"], sub))
	if err != nil {
		return nil, err
	}
	rresp := new(redditResponse)
	err = json.Unmarshal(b, rresp)
	rresp.Data.Name = sub
	return &rresp.Data, nil
}

//GetSubredditN gets the first n items from a subreddit
//returns a subreddit object containing the items
func GetSubredditN(sub string, n int) (*Subreddit, error) {
	var tempsub *Subreddit
	for ; n > 0; n -= 100 {
		data := url.Values{
			"limit": {fmt.Sprintf("%d", n)},
		}
		if tempsub != nil && len(tempsub.Items) != 0 {
			data.Set("after", tempsub.Items[len(tempsub.Items)-1].Name)
		}
		b, err := getPostJsonBytes(fmt.Sprintf(Urls["subreddit"], sub), &data)
		if err != nil {
			return nil, err
		}
		rresp := new(redditResponse)
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
func GetFrontPage(user *Redditor) (*Subreddit, error) {
	b, err := getJsonBytes(Urls["frontpage"])
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
	b, err := getJsonBytes(fmt.Sprintf(Urls["user"], user))
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
	data := url.Values{
		"user":     {user},
		"passwd":   {pass},
		"api_type": {"json"},
		"rem":      {fmt.Sprintf("%v", rem)},
	}
	b, err := getPostJsonBytes(ApiUrls["login"], &data)
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
