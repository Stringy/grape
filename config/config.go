package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var conf = new(cfg)

type cfg struct {
	UserAgent string `json: "user_agent"`
	Host      string
	ApiUrls   map[string]string `json: "api_urls"`
	Urls      map[string]string
}

func init() {
	err := Load("config.json")
	if err != nil {
		panic(err)
	}
}

func Load(fn string) error {
	f, err := os.Open("config.json")
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

func ApiUrl(name string) string {
	return conf.Host + conf.ApiUrls[name]
}

func Url(name string) string {
	return conf.Host + conf.Urls[name]
}

func Host() string {
	return conf.Host
}

func UserAgent() string {
	return conf.UserAgent
}
