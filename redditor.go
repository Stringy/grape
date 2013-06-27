package reddit

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

type Redditor struct {
	Name     string
	LKarma   int  `json:"link_karma"`
	CKarma   int  `json:"comment_karma"`
	IsFriend bool `json:"is_friend"`
	HasMail  bool `json:"has_mail"`
	IsOver18 bool `json:"over_18"`
	IsGold   bool `json:"is_gold"`
	IsMod    bool `json:"is_mod"`
	ModHash  string
}

func (r *Redditor) IsLoggedIn() bool {
	return len(r.ModHash) != 0
}

func (r *Redditor) MakeComment(parent, body string) error {
	return nil
}

func (r *Redditor) SubmitLink(subreddit, title, body, link, kind string) error {
	if r == nil {
		return errors.New("reddit: nil redditor")
	}
	if !r.IsLoggedIn() {
		return errors.New("Submission error: User is not logged in")
	}
	l, _ := url.Parse(submit_url)
	data := url.Values{
		"api_type":    {"json"},
		"captcha":     {"this is a captcha"},
		"resubmit":    {"true"},
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

	respBytes, err := getPostJsonBytes(l, &data)
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

func (r *Redditor) DeleteAccount(passwd string) error {
	if r == nil || !r.IsLoggedIn() {
		return errors.New("reddit: can't delete redditor without logging in")
	}
	link, _ := url.Parse(delete_url)
	data := url.Values{
		"api_type":       {"json"},
		"confirm":        {"false"},
		"delete_message": {""},
		"passwd":         {passwd},
		"uh":             {r.ModHash},
		"user":           {r.Name},
	}
	respBytes, err := getPostJsonBytes(link, &data)
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

func (r *Redditor) GetUnreadMail() ([]string, error) {
	return nil, nil
}

func (r *Redditor) GetFrontpage() (*Subreddit, error) {
	return nil, nil
}

func (r *Redditor) Me() error {
	if !r.IsLoggedIn() {
		return errors.New("reddit: user not logged in")
	}
	link, _ := url.Parse(me_url)
	respBytes, err := getJsonBytes(link)
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
