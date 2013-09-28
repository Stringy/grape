package reddit

import (
	"errors"
	"io/ioutil"
	logging "log"
	"os"
)

var log = logging.New(ioutil.Discard, "[reddit] ", logging.LstdFlags)
var debug = logging.New(ioutil.Discard, "[reddit debug] ", logging.LstdFlags)

func init() {
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
		debug = logging.New(out, "[reddit debug] ", logging.LstdFlags)
	}
}

// Reusable Errors
var (
	NotLoggedInError  = errors.New("reddit: user not logged in")
	TitleTooLongError = errors.New("reddit: title too long; must be <= 300 characters")
)

// Error json response
type errorJson struct {
	Json struct {
		Errors [][]string
	}
}
