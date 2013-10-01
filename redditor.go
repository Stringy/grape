package grape

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Redditor struct {
	LKarma        int  `json:"link_karma"`
	CKarma        int  `json:"comment_karma"`
	IsFriend      bool `json:"is_friend"`
	HasMail       bool `json:"has_mail"`
	IsOver18      bool `json:"over_18"`
	IsGold        bool `json:"is_gold"`
	IsMod         bool `json:"is_mod"`
	Created       int
	CreatedUTC    time.Time `json:"created_utc"`
	HasModMail    bool      `json:"has_mod_mail"`
	VerifiedEmail bool      `json:"has_verified_email"`
	ModHash       string
	*Thing
}

func NewRedditor() *Redditor {
	r := new(Redditor)
	r.Thing = new(Thing)
	return r
}

//IsLoggedIn returns true if the user is currently logged into reddit
func (r *Redditor) IsLoggedIn() bool {
	if len(r.ModHash) == 0 || len(client.Jar.Cookies(reddit_url)) == 0 {
		return false
	}
	for _, cookie := range client.Jar.Cookies(reddit_url) {
		if cookie.Expires.After(time.Now()) {
			return true
		} else {
			debug.Println("Expired cookie: ", cookie)
		}
	}
	return false
}

//ReplyToComment replies to a reddit comment on behalf of the user
func (r *Redditor) ReplyToComment(parent *Comment, body string) error {
	return parent.Reply(r, body)
}

//PostComment posts a top level comment to a reddit submission
func (r *Redditor) PostComment(parent *Submission, body string) error {
	return parent.PostComment(r, body)
}

//SubmitSelf submits a self (non-link) submission to the subreddit of choice
func (r *Redditor) SubmitSelf(subreddit, title, body string) error {
	return r.submit(subreddit, title, body, "", "self", true)
}

//SubmitLink submits a link to the subreddit
func (r *Redditor) SubmitLink(subreddit, title, link string, resubmit bool) error {
	return r.submit(subreddit, title, "", link, "link", resubmit)
}

//submit handles all the submission semantics for the top level functions
func (r *Redditor) submit(subreddit, title, body, link, kind string, resubmit bool) error {
	if r == nil {
		return errors.New("nil redditor")
	}
	if len(title) > 300 {
		return titleTooLongError
	}
	if !r.IsLoggedIn() {
		return notLoggedInError
	}

	data := url.Values{
		"api_type":    {"json"},
		"captcha":     {""},
		"resubmit":    {fmt.Sprintf("%v", resubmit)},
		"extension":   {"/"},
		"iden":        {""},
		"save":        {"false"},
		"kind":        {kind},
		"then":        {"comments"},
		"sr":          {subreddit},
		"text":        {body},
		"title":       {title},
		"uh":          {r.ModHash},
		"url":         {link},
		"sendreplies": {"true"},
	}

	errstruct := new(errorJson)
	respBytes, err := makePostRequest(Config.GetApiUrl("submit"), &data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBytes, &errstruct)
	if err != nil {
		return err
	}
	if len(errstruct.Json.Errors) != 0 {
		return errors.New(strings.Join(errstruct.Json.Errors[0], ", "))
	}
	return nil
}

//DeleteAccount deletes the user from reddit
func (r *Redditor) DeleteAccount(passwd string) error {
	if r == nil || !r.IsLoggedIn() {
		return notLoggedInError
	}
	data := url.Values{
		"api_type":       {"json"},
		"confirm":        {"false"},
		"delete_message": {""},
		"passwd":         {passwd},
		"uh":             {r.ModHash},
		"user":           {r.Name},
	}
	respBytes, err := makePostRequest(Config.GetApiUrl("delete"), &data)
	if err != nil {
		return err
	}
	errstruct := new(struct {
		Json struct {
			Errors [][]string
		}
	})
	err = json.Unmarshal(respBytes, errstruct)
	if err != nil {
		return err
	}
	if len(errstruct.Json.Errors) != 0 {
		return errors.New(strings.Join(errstruct.Json.Errors[0], ", "))
	}
	return nil
}

// GetUnreadMail gets the unread mail for the user
// doesn't require modhash for reading messages
func (r *Redditor) GetUnreadMail(limit int) ([]Message, error) {
	data := url.Values{
		"limit": {fmt.Sprintf("%d", limit)},
	}
	b, err := makeGetRequest(Config.GetUrl("unread"), &data)
	if err != nil {
		return nil, err
	}
	//debug.Println("unread messages json:", string(b))
	msgresp := new(messageResponse)
	err = json.Unmarshal(b, msgresp)
	if err != nil {
		return nil, err
	}
	msgs := make([]Message, len(msgresp.Data.Children))
	for i, msg := range msgresp.Data.Children {
		msgs[i] = msg.Msg
	}
	return msgs, nil
}

// GetInbox gets all mail from the user's mail
// doesn't require modhash for reading
func (r *Redditor) GetInbox(limit int) ([]Message, error) {
	data := url.Values{
		"limit": {fmt.Sprintf("%d", limit)},
	}
	b, err := makeGetRequest(Config.GetUrl("inbox"), &data)
	if err != nil {
		return nil, err
	}
	msgresp := new(messageResponse)
	err = json.Unmarshal(b, msgresp)
	if err != nil {
		return nil, err
	}
	debug.Println(msgresp.Data.Children)
	msgs := make([]Message, len(msgresp.Data.Children))
	for i, msg := range msgresp.Data.Children {
		msgs[i] = msg.Msg
	}
	return msgs, nil
}

// GetFrontpage returns the frontpage for the user, including all
// subscribed subreddits
func (r *Redditor) GetFrontpage() (*Subreddit, error) {
	return nil, nil
}

// GetFrontpageN returns the first n submissions from their frontpage
// similar semantics to GetSubredditN
func (r *Redditor) GetFrontpageN(n int) (*Subreddit, error) {
	return nil, nil
}

// Me populates the redditor with their information
func (r *Redditor) Me() error {
	if !r.IsLoggedIn() {
		return notLoggedInError
	}
	respBytes, err := makeGetRequest(Config.GetApiUrl("me"), nil)
	if err != nil {
		return err
	}
	uresp := new(userResponse)
	uresp.Data = *r
	err = json.Unmarshal(respBytes, uresp)
	if err != nil {
		return nil
	}
	r = &uresp.Data
	return nil
}

// ClearSessions clears the current user's reddit sessions
// This will result in immediate logging out of the user
func (r *Redditor) ClearSessions(pass string) error {
	data := url.Values{
		"api_type": {"json"},
		"curpass":  {pass},
		"uh":       {r.ModHash},
		"dest":     {Config.Host},
	}
	b, err := makePostRequest(Config.GetApiUrl("clear_sessions"), &data)
	if err != nil {
		return err
	}
	//debug.Println(string(b))
	erresp := new(errorJson)
	err = json.Unmarshal(b, erresp)
	if err != nil {
		return err
	}
	if len(erresp.Json.Errors) != 0 {
		return errors.New(strings.Join(erresp.Json.Errors[0], ", "))
	}
	return nil
}
