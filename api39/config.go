package api39

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Filename string `json:"-"`
	Apikey   string `json:"apikey"`
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

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(Config{"", apikey})
}
