package reddit

import (
	"errors"
)

var (
	NotLoggedInError = errors.New("reddit: user not logged in")
)
