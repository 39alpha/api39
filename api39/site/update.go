package site

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/39alpha/api39/api39"
	"github.com/kataras/iris/v12"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	godaddyapi = "https://api.godaddy.com/v1/domains/"
	mainbranch = "origin/main"
	devbranch  = "origin/dev"
)

func UpdateDNSLink(cfg *api39.Config, ipfshash string) error {
	url := godaddyapi + cfg.Domain + "/records/TXT/_dnslink"

	payload := []map[string]string{
		map[string]string{"data": "dnslink=/ipfs/" + ipfshash},
	}
	content, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	data := bytes.NewReader(content)

	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest(http.MethodPut, url, data)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "sso-key "+cfg.GoDaddy.Key+":"+cfg.GoDaddy.Secret)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode != 200 {
		return fmt.Errorf("%s", res.Status)
	}

	return nil
}

func Update(ctx iris.Context) {
	body, ok := ctx.Values().Get("JSONBody").(iris.Map)
	if !ok {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("Bad request body"))
	}

	if ref, ok := body["ref"]; !ok {
		log.Println("Request Recieved: no reference key")
		_, _ = ctx.JSON(iris.Map{"message": "successful request"})
	} else if ref.(string) == "refs/heads/main" {
		go UpdateProduction(ctx)
		_, _ = ctx.JSON(iris.Map{"message": "successful"})
	} else if ref.(string) == "refs/heads/dev" {
		go UpdateDevelopment(ctx)
		_, _ = ctx.JSON(iris.Map{"message": "successful"})
	} else {
		log.Printf("Request Recieved: branch is neither main nor dev (%v)\n", ref.(string))
		_, _ = ctx.JSON(iris.Map{"message": "successful request"})
	}
}

func UpdateProduction(ctx iris.Context) {
	cfg, ok := ctx.Values().Get("config").(*api39.Config)
	if !ok {
		log.Println("failed to retrieve configuration from context")
		return
	}

	if err := api39.UpdateGitRepo(cfg.Site.Repo, cfg.Site.Path, mainbranch); err != nil {
		message := "failed to update repository"
		log.Printf("%s: %v\n", message, err)
		return
	}

	if err := api39.RebuildWithMake(cfg.Site.Path); err != nil {
		message := "failed to rebuild site"
		log.Printf("%s: %v\n", message, err)
		return
	}

	path := filepath.Join(cfg.Site.Path, "public")
	hash, err := api39.IpfsAddDir(cfg.Ipfs.Url, path)
	if err != nil {
		message := "failed to add site to IPFS"
		log.Printf("%s: %v\n", message, err)
		return
	}

	log.Printf("New IPFS Hash: %s\n", hash)

	err = UpdateDNSLink(cfg, hash)
	if err != nil {
		message := "failed to update DNS records"
		log.Printf("%s: %v\n", message, err)
		return
	}
}

func UpdateDevelopment(ctx iris.Context) {
	cfg, ok := ctx.Values().Get("config").(*api39.Config)
	if !ok {
		log.Println("failed to retrieve configuration from context")
		return
	}

	if err := api39.UpdateGitRepo(cfg.Site.Repo, cfg.Site.DevPath, devbranch); err != nil {
		message := "failed to update repository"
		log.Printf("%s: %v\n", message, err)
		return
	}

	if err := api39.RebuildWithMake(cfg.Site.DevPath, "dev"); err != nil {
		message := "failed to rebuild site"
		log.Printf("%s: %v\n", message, err)
		return
	}
}
