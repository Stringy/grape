package reddit

import (
	"net/http"
	"net/url"
)

var (
	// requests
	get       *http.Request
	post      *http.Request
	post_form *http.Request
	client    *http.Client

	actual_url *url.URL

	UserAgent = "Go Reddit API by String217 v0.1"

	ApiUrls = map[string]string{
		"login":          host + "/api/login",
		"me":             host + "/api/me.json",
		"comment":        host + "/api/comment",
		"delete_user":    host + "/api/delete_user",
		"captcha":        host + "/api/new_captcha",
		"submit":         host + "/api/submit",
		"user_avail":     host + "/api/username_available.json",
		"clear_sessions": host + "/api/clear_sessions",
		"register":       host + "/api/register",
		"update":         host + "/api/update",
		"del":            host + "/api/del",
		"editusertext":   host + "/api/editusertext",
		"hide":           host + "/api/hide",
		"info":           host + "/api/info",
		"marknsfw":       host + "/api/marknsfw",
		"morechildren":   host + "/api/morechildren",
		"report":         host + "/api/report",
		"save":           host + "/api/save",
		"unhide":         host + "/api/unhide",
		"unmarknsfw":     host + "/api/unmarknsfw",
		"vote":           host + "/api/vote",
		"block":          host + "/api/block",
		"compose":        host + "/api/compose",
		"read_message":   host + "/api/read_message",
		"unread_message": host + "/api/unread_message",
	}

	Urls = map[string]string{
		"home":      host,
		"subreddit": "/r/%s.json",
		"frontpage": "/.json",
		"user":      "/user/%s/about.json",
		"comment":   "/r/%s/%s.json",
	}
)

const (
	host = "http://reddit.local"

	config_file = "reddit.conf"

	KindLink = "link"
	KindSelf = "self"
)

func init() {
	actual_url, _ = url.Parse("http://reddit.local/")
	jar := NewJar()
	client = &http.Client{nil, nil, jar}
	var err error
	get, err = http.NewRequest("GET", "", nil)
	if err != nil {
		panic(err)
	}
	post, err = http.NewRequest("POST", "", nil)
	if err != nil {
		panic(err)
	}
	post_form, err = http.NewRequest("POST", "", nil)
	if err != nil {
		panic(err)
	}
	post_form.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	get.Header.Set("User-Agent", UserAgent)
	post.Header.Set("User-Agent", UserAgent)
	post_form.Header.Set("User-Agent", UserAgent)
}
