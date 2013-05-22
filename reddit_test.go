package reddit

import (
	"fmt"
	"testing"
)

func TestReddit(t *testing.T) {
	r := GetFrontPage()
	//fmt.Println(r)
	i := r.Items[0]
	fmt.Println(i.Ups, i.Downs, i.Score)
}
