package reddit

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var config = new(cfg)

type cfg struct {
	UserAgent string            `json:"user_agent"`
	Host      string            `json:"host"`
	ApiUrl    map[string]string `json:"api_urls"`
	Url       map[string]string `json:"urls"`
	Log       bool              `json:"enable_logging"`
	Debug     bool              `json:"enable_debug"`
	DebugFile string            `json:"debug_file"`
	LogFile   string            `json:"log_file"`
}

func init() {
	err := Load("config.json")
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
	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *cfg) FullApiUrl(name string) string {
	return c.Host + c.ApiUrl[name]
}

func (c *cfg) FullUrl(name string) string {
	return c.Host + c.Url[name]
}
