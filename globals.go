package grape

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	logging "log"
	"net/url"
	"os"
	"runtime"
)

var reddit_url *url.URL

// Logging
var log = logging.New(ioutil.Discard, "[reddit] ", logging.LstdFlags)
var debug = logging.New(ioutil.Discard, "[reddit debug] ", logging.LstdFlags)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}
	if config.Log {
		out, err := os.Create(config.LogFile)
		if err != nil {
			panic(err)
		}
		log = logging.New(out, "[reddit] ", logging.LstdFlags)
	}
	if config.Debug {
		out, err := os.Create(config.DebugFile)
		if err != nil {
			panic(err)
		}
		debug = logging.New(out, "[reddit debug] ", logging.LstdFlags|logging.Lshortfile)
	}
	reddit_url, err = url.Parse(config.Host)
	if err != nil {
		log.Panicf("error parsing config.Host into url: %v", err)
	}
}

// Configuration structures
var config = new(cfg)

type cfg struct {
	UserAgent string            `json:"user_agent"`
	Host      string            `json:"host"`
	ApiUrl    map[string]string `json:"api_urls"`
	Url       map[string]string `json:"urls"`
	Log       bool              `json:"enable_logging"`
	Debug     bool              `json:"enable_debug"`
	DebugFile string            `json:"debug_file"`
	LogFile   string            `json:"log_file"`
}

// loadConfig decodes the configuration information from the config file
func loadConfig(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}
	return nil
}

// GetApiUrl gives the api url including host
func (c *cfg) GetApiUrl(name string) string {
	return c.Host + c.ApiUrl[name]
}

// GetUrl gives the reddit url format string including host
func (c *cfg) GetUrl(name string) string {
	return c.Host + c.Url[name]
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
