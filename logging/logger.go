package logging

import (
	"log"
	"os"
	"reddit/config"
)

var rlog *log.Logger
var cfg = config.GetInstance()

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

func Print(v ...interface{}) {
	if cfg.Log {
		rlog.Print(v...)
	}
}

func Printf(f string, v ...interface{}) {
	if cfg.Log {
		rlog.Printf(f, v...)
	}
}

func Println(v ...interface{}) {
	if cfg.Log {
		rlog.Println(v...)
	}
}

func Fatalf(f string, v ...interface{}) {
	if cfg.Log {
		rlog.Fatalf(f, v...)
	}
}
