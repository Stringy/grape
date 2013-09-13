package logging

import (
	"log"
	"os"
)

var rlog *log.Logger

const (
	prefix = "[reddit]"
)

func init() {
	f, err := os.Create("reddit.log")
	if err != nil {
		panic(err)
	}
	rlog = log.New(f, prefix, log.LstdFlags)
}
