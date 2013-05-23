package reddit

import (
	"fmt"
	"testing"
)

func TestReddit(t *testing.T) {
	_, err := Login("", "", false)
	if err != nil {
		fmt.Println(err)
	}
}
