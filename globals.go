package grape

import (
	"errors"
	"io/ioutil"
	logging "log"
	"net/url"
	"os"
	"runtime"
)

var reddit_url *url.URL

// Standard Logging
var log = logging.New(ioutil.Discard, "[reddit] ", logging.LstdFlags)

// Debugging
var debug = logging.New(ioutil.Discard, "[reddit debug] ", logging.LstdFlags)

// sorting type
type sort string

// Listing sort constants
const (
	HOT       sort = "hot"
	TOP       sort = "top"
	NEW       sort = "new"
	CONT      sort = "controversial"
	DEFAULT_S sort = "subreddit"
)

// time period for sorting
type period string

// time period constants
const (
	HOUR      period = "hour"
	DAY       period = "day"
	WEEK      period = "week"
	MONTH     period = "month"
	YEAR      period = "year"
	ALL       period = "all"
	DEFAULT_P period = ""
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initConfig()
	if Config.Log {
		out, err := os.Create(Config.LogFile)
		if err != nil {
			panic(err)
		}
		log = logging.New(out, "[reddit] ", logging.LstdFlags)
	}
	if Config.Debug {
		out, err := os.Create(Config.DebugFile)
		if err != nil {
			panic(err)
		}
		debug = logging.New(out, "[reddit debug] ", logging.LstdFlags|logging.Lshortfile)
	}
	var err error
	reddit_url, err = url.Parse(Config.Host)
	if err != nil {
		log.Panicf("error parsing Config.Host into url: %v", err)
	}
}

// Configuration structures
var Config = new(cfg)

type cfg struct {
	UserAgent string
	Host      string
	apiUrl    map[string]string
	url       map[string]string
	Log       bool
	Debug     bool
	DebugFile string
	LogFile   string
}

// initConfig decodes the Configuration information from the Config file
func initConfig() {
	Config.UserAgent = "/u/stringy217's Go reddit api v0.1"
	Config.Host = "http://reddit.local"
	Config.apiUrl = map[string]string{
		"login":          "/api/login",
		"me":             "/api/me.json",
		"comment":        "/api/comment",
		"delete_user":    "/api/delete_user",
		"captcha":        "/api/new_captcha",
		"submit":         "/api/submit",
		"user_avail":     "/api/username_available.json",
		"clear_sessions": "/api/clear_sessions",
		"register":       "/api/register",
		"update":         "/api/update",
		"del":            "/api/del",
		"editusertext":   "/api/editusertext",
		"hide":           "/api/hide",
		"info":           "/api/info",
		"marknsfw":       "/api/marknsfw",
		"morechildren":   "/api/morechildren",
		"report":         "/api/report",
		"save":           "/api/save",
		"unhide":         "/api/unhide",
		"unmarknsfw":     "/api/unmarknsfw",
		"vote":           "/api/vote",
		"block":          "/api/block",
		"compose":        "/api/compose",
		"read_message":   "/api/read_message",
		"unread_message": "/api/unread_message",
	}
	Config.url = map[string]string{
		"subreddit":     "/r/%s.json",
		"limited_sub":   "/r/%s/",
		"frontpage":     "/.json",
		"user":          "/user/%s/about.json",
		"comment":       "/r/%s/%s.json",
		"inbox":         "/message/inbox.json",
		"unread":        "/message/unread.json",
		"sent":          "/message/sent.json",
		"hot":           "/r/%s/hot.json",
		"new":           "/r/%s/new.json",
		"controversial": "/r/%s/controversial.json",
	}
	Config.Log = true
	Config.Debug = true
	Config.LogFile = "reddit.log"
	Config.DebugFile = "reddit.debug.log"
}

// GetApiUrl gives the api url including host
func (c *cfg) GetApiUrl(name string) string {
	return c.Host + c.apiUrl[name]
}

// GetUrl gives the reddit url format string including host
func (c *cfg) GetUrl(name string) string {
	return c.Host + c.url[name]
}

func (c *cfg) SetUserAgent(ua string) {
	c.UserAgent = ua
}

// Reusable Errors
var (
	notLoggedInError    = errors.New("reddit: user not logged in")
	titleTooLongError   = errors.New("reddit: title too long; must be <= 300 characters")
	incorrectOwnerError = errors.New("reddit: user does not have ownership over reddit thing")
)

// Error json response
type errorJson struct {
	Json struct {
		Errors [][]string
	}
}
