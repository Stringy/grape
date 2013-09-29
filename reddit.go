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
	b, err := makeGetRequest(fmt.Sprintf(Config.GetUrl("subreddit"), sub))
	if err != nil {
		return nil, err
	}
	rresp := new(redditResponse)
	err = json.Unmarshal(b, rresp)
	rresp.Data.Name = sub
	return &rresp.Data, nil
}

// GetSortedSubreddit gets the front page of a named subreddit with submissions sorted.
//
// sub is the name of a valid subreddit.
//
// s is the required sorting of the subreddit (hot, new, controversial, top or default).
//
// p is the required time period for the sorting (hour, day, week, month, year, all).
func GetSortedSubreddit(sub string, s sort, p period) (*Subreddit, error) {
	log.Printf("getting sorted subreddit %s sorted by %s over period %s", sub, s, p)
	data := url.Values{
		"t": {string(p)},
	}
	b, err := makePostRequest(fmt.Sprintf(Config.GetUrl(string(s)), sub), &data)
	if err != nil {
		return nil, err
	}
	rresp := new(redditResponse)
	err = json.Unmarshal(b, rresp)
	rresp.Data.Name = sub
	return &rresp.Data, nil
}

// GetSubredditN gets the first n items from a subreddit.
//
// returns a subreddit object containing the items.
//
// sub is the name of a valid subreddit.
func GetSubredditN(sub string, n int) (*Subreddit, error) {
	u := fmt.Sprintf(Config.GetUrl("limited_sub"), sub)
	var tempsub *Subreddit
	for ; n > 0; n -= 100 {
		data := url.Values{
			"limit":    {fmt.Sprintf("%d", n)},
			"api_type": {"json"},
		}
		if tempsub != nil && len(tempsub.Items) != 0 {
			data.Set("after", tempsub.Items[len(tempsub.Items)-1].Name)
		}
		b, err := makePostRequest(u, &data)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(b))
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

// GetSubredditN gets the first n items from a subreddit with submissions sorted.
//
// returns a subreddit object containing the items
//
// sub is the name of a valid subreddit
//
// s is the required sorting of the subreddit (hot, new, controversial, top or default) can be nil
//
// p is the required time period for the sorting (hour, day, week, month, year, all) can be nil
//
// n is the number of required submissions
func GetSortedSubredditN(sub string, s sort, p period, n int) (*Subreddit, error) {
	return nil, nil
}

// GetFrontPage currently gets the front page of *default* reddit
// TODO(Stringy): apply this to currently logged in user
func GetFrontPage(user *Redditor) (*Subreddit, error) {
	b, err := makeGetRequest(Config.GetUrl("frontpage"))
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
	b, err := makeGetRequest(fmt.Sprintf(Config.GetUrl("user"), user))
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
