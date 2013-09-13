package things

import (
	"errors"
)

var (
	NotLoggedInError  = errors.New("reddit: user not logged in")
	TitleTooLongError = errors.New("reddit: title too long; must be <= 300 characters")
)

type errorJson struct {
	Json struct {
		Errors [][]string
	}
}
