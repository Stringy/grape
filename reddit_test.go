package reddit

import (
	"fmt"
	"testing"
)

func TestReddit(t *testing.T) {
	r := GetSubreddit("learnprogramming")
	post := r.Items[0]
	p := post.GetComments()
	fmt.Println(len(p))
	fmt.Println(p[0])
}
