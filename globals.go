package reddit

import (
	"net/http"
	"net/url"
)

var (
	// requests
	client *http.Client

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
		"subreddit": host + "/r/%s.json",
		"frontpage": host + "/.json",
		"user":      host + "/user/%s/about.json",
		"comment":   host + "/r/%s/%s.json",
	}
)

const (
	host = "http://reddit.local"

	config_file = "reddit.conf"
)

func init() {
	actual_url, _ = url.Parse("http://reddit.local/")
	jar := NewJar()
	client = &http.Client{nil, nil, jar}
}
