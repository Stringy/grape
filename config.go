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
	apiUrl    map[string]string `json:"api_urls"`
	url       map[string]string `json:"urls"`
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

func (c *cfg) GetApiUrl(name string) string {
	return c.Host + c.apiUrl[name]
}

func (c *cfg) GetUrl(name string) string {
	return c.Host + c.url[name]
}
