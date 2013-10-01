package grape

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// GetSubreddit gets the front page of a named subreddit.
//
// sub is the name of a valid subreddit
func GetSubreddit(sub string) (*Subreddit, error) {
	log.Printf("Getting subreddit: %s\n", sub)
	b, err := makeGetRequest(fmt.Sprintf(Config.GetUrl("subreddit"), sub), nil)
	if err != nil {
		return nil, err
	}
	rresp := new(redditResponse)
	err = json.Unmarshal(b, rresp)
	if err != nil {
		return nil, err
	}
	rresp.Data.Name = sub
	return &rresp.Data, nil
}

// GetFrontPage currently gets the front page of *default* reddit
// TODO(Stringy): apply this to currently logged in user
func GetFrontPage(user *Redditor) (*Subreddit, error) {
	b, err := makeGetRequest(Config.GetUrl("frontpage"), nil)
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

// GetRedditor returns a redditor object containing all information relevant to that
// reddit user
func GetRedditor(user string) (*Redditor, error) {
	log.Printf("getting Redditor: %s\n", user)
	b, err := makeGetRequest(fmt.Sprintf(Config.GetUrl("user"), user), nil)
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

// Login logs a user into reddit through the api login page
// returns the same errors recieved from reddit, if applicable
// otherwise returns a redditor with populated modhash
// TODO(Stringy): add support for ssl.reddit.com/login
func Login(user, pass string, rem bool) (*Redditor, error) {
	log.Printf("logging in to user: %s\n", user)
	data := url.Values{
		"user":     {user},
		"passwd":   {pass},
		"api_type": {"json"},
		"rem":      {fmt.Sprintf("%v", rem)},
	}
	// debug.Println("making login request")
	b, err := makePostRequest(Config.GetApiUrl("login"), &data)
	if err != nil {
		return nil, err
	}
	// debug.Println("response body:", string(b))
	loginResp := new(loginResponse)
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
