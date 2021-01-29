package api39

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

type StripeConfig struct {
	Apikey     string `json:"apikey"`
	Currency   string `json:"currency"`
	Product    string `json:"product"`
	SuccessURL string `json:"success"`
	CancelURL  string `json:"cancel"`
}

type SiteConfig struct {
	Repo string `json:"repo"`
	Path string `json:"path"`
	Hugo string `json:"hugo"`
}

type IpfsConfig struct {
	Url string `json:"url"`
}

type Config struct {
	Filename string       `json:"-"`
	Apikey   string       `json:"apikey"`
	Site     SiteConfig   `json:"site"`
	Ipfs     IpfsConfig   `json:"ipfs"`
	Stripe   StripeConfig `json:"stripe"`
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

	hugopath, err := exec.LookPath("hugo")
	if err != nil {
		hugopath = ""
	}

	site := SiteConfig{
		"https://github.com/39alpha/39alpharesearch.org",
		"",
		hugopath,
	}

	stripe := StripeConfig{
		"",
		"usd",
		"Your Generous Donation",
		"https://39alpharesearch.org/donate/success",
		"https://39alpharesearch.org/donate",
	}

	ipfs := IpfsConfig{"127.0.0.1:5001"}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(Config{"", apikey, site, ipfs, stripe})
}
