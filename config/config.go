package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var conf = new(cfg)

type cfg struct {
	UserAgent string            `json:"user_agent"`
	Host      string            `json:"host"`
	ApiUrl    map[string]string `json:"api_urls"`
	Url       map[string]string `json:"urls"`
	Log       bool              `json:"enable_logging"`
	LogFile   string            `json:"log_file"`
}

func init() {
	err := Load("../config.json")
	if err != nil {
		panic(err)
	}
}

func Load(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return err
	}
	return nil
}

func GetInstance() *cfg {
	return conf
}

// func ApiUrl(name string) string {
// 	return conf.Host + conf.ApiUrls[name]
// }

// func Url(name string) string {
// 	return conf.Host + conf.Urls[name]
// }

// func Host() string {
// 	return conf.Host
// }

// func UserAgent() string {
// 	return conf.UserAgent
// }

// func Log() bool {
// 	return conf.Log
// }
