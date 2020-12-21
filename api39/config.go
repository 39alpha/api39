package api39

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type SiteConfig struct {
	Repo string `json:"repo"`
	Path string `json:"path"`
}

type Config struct {
	Filename string     `json:"-"`
	Apikey   string     `json:"apikey"`
	Site     SiteConfig `json:"site"`
}

func ReadConfig(filename string) (*Config, error) {
	blob, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(blob, &config); err != nil {
		return nil, err
	}
	config.Filename = filename

	return &config, nil
}

func GenerateConfig(n int) error {
	apikey, err := GenerateApiKey(n)
	if err != nil {
		return err
	}

	site := SiteConfig{
		"https://github.com/39alpha/39alpharesearch.org",
		"",
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(Config{"", apikey, site})
}
