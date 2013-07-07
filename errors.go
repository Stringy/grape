package reddit

import (
	"errors"
)

var (
	TitleTooLongError = errors.New("reddit: title too long; must be <= 300 characters")
	NotLoggedInError  = errors.New("reddit: user not logged in")
)

type errorJson struct {
	Json struct {
		Errors [][]string
	}
}
