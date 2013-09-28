package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// GetSubreddit gets the front page of a named subreddit
// TODO: add support for arbitrary number of posts returned
func GetSubreddit(sub string) (*Subreddit, error) {
	log.Printf("Getting subreddit: %s\n", sub)
	b, err := makeGetRequest(fmt.Sprintf(config.GetUrl("subreddit"), sub))
	if err != nil {
		return nil, err
	}
	rresp := new(RedditResponse)
	err = json.Unmarshal(b, rresp)
	rresp.Data.Name = sub
	return &rresp.Data, nil
}

// GetSubredditN gets the first n items from a subreddit
// returns a subreddit object containing the items
func GetSubredditN(sub string, n int) (*Subreddit, error) {
	var tempsub *Subreddit
	for ; n > 0; n -= 100 {
		data := url.Values{
			"limit": {fmt.Sprintf("%d", n)},
		}
		if tempsub != nil && len(tempsub.Items) != 0 {
			data.Set("after", tempsub.Items[len(tempsub.Items)-1].Name)
		}
		b, err := makePostRequest(fmt.Sprintf(config.GetUrl("subreddit"), sub), &data)
		if err != nil {
			return nil, err
		}
		rresp := new(RedditResponse)
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
	b, err := makeGetRequest(config.GetApiUrl("frontpage"))
	if err != nil {
		return nil, err
	}
	rresp := new(RedditResponse)
	err = json.Unmarshal(b, rresp)
	if err != nil {
		return nil, err
	}
	return &rresp.Data, nil
}

// GetRedditor returns information about a given redditor
func GetRedditor(user string) (*Redditor, error) {
	log.Printf("getting Redditor: %s\n", user)
	b, err := makeGetRequest(fmt.Sprintf(config.GetUrl("user"), user))
	if err != nil {
		return nil, err
	}
	uresp := new(UserResponse)
	err = json.Unmarshal(b, uresp)
	if err != nil {
		return nil, err
	}
	return &uresp.Data, nil
}

// Login logs a user into reddit through the api login page
// returns the same errors recieved from reddit, if applicable
// otherwise returns a redditor with populated modhash
func Login(user, pass string, rem bool) (*Redditor, error) {
	log.Printf("logging in to user: %s\n", user)
	data := url.Values{
		"user":     {user},
		"passwd":   {pass},
		"api_type": {"json"},
		"rem":      {fmt.Sprintf("%v", rem)},
	}
	// debug.Println("making login request")
	b, err := makePostRequest(config.GetApiUrl("login"), &data)
	if err != nil {
		return nil, err
	}
	// debug.Println("response body:", string(b))
	loginResp := new(LoginResponse)
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

	redditor := NewRedditor()
	redditor.Name = user
	redditor.ModHash = loginResp.Json.Data.ModHash
	return redditor, nil
}
