package reddit

import (
	"fmt"
	"testing"
)

func TestReddit(t *testing.T) {
	r := GetFrontPage()
	//fmt.Println(r)
	for _, i := range r.Items {
		fmt.Println(i.String())
	}
}
