package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Redditor struct {
	LKarma   int  `json:"link_karma"`
	CKarma   int  `json:"comment_karma"`
	IsFriend bool `json:"is_friend"`
	HasMail  bool `json:"has_mail"`
	IsOver18 bool `json:"over_18"`
	IsGold   bool `json:"is_gold"`
	IsMod    bool `json:"is_mod"`
	ModHash  string
	*Thing
}

//IsLoggedIn returns true if the user is currently logged into reddit
func (r *Redditor) IsLoggedIn() bool {
	return len(r.ModHash) != 0
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
		return errors.New("reddit: nil redditor")
	}
	if len(title) > 300 {
		return TitleTooLongError
	}
	if !r.IsLoggedIn() {
		return errors.New("Submission error: User is not logged in")
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

	errstruct := new(struct {
		Json struct {
			Errors [][]string
		}
	})

	respBytes, err := getPostJsonBytes(ApiUrls["submit"], &data)
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
		return errors.New("reddit: can't delete redditor without logging in")
	}
	data := url.Values{
		"api_type":       {"json"},
		"confirm":        {"false"},
		"delete_message": {""},
		"passwd":         {passwd},
		"uh":             {r.ModHash},
		"user":           {r.Name},
	}
	respBytes, err := getPostJsonBytes(ApiUrls["delete"], &data)
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

//GetUnreadMail gets the unread mail for the user
func (r *Redditor) GetUnreadMail() ([]string, error) {
	return nil, nil
}

//GetFrontpage returns the frontpage for the user, including all 
//subscribed subreddits
func (r *Redditor) GetFrontpage() (*Subreddit, error) {
	return nil, nil
}

//Me populates the redditor with their information
func (r *Redditor) Me() error {
	if !r.IsLoggedIn() {
		return errors.New("reddit: user not logged in")
	}
	respBytes, err := getJsonBytes(ApiUrls["me"])
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
