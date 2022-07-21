package site

import (
	"fmt"
	"github.com/39alpha/api39/api39"
	"github.com/kataras/iris/v12"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"encoding/json"
)

const (
	godaddyapi := "https://api.godaddy.com/v1/domains/"
)

func UpdateDNSLink(cfg api39.Config, ipfshash string) error {
	url := godaddyapi + cfg.Domain + "/records/TXT/_dnslink"

	payload := []map[string]string{
		map[string]string{ "data": "/ipfs/" + ipfshash },
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
	req.Header.Add("Authorization", "sso-key " + key + ":" + secret)
	fmt.Println(req.Header)
	if err != nil {
		return err;
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	} else if res.StatusCode != 200 {
		return fmt.Errorf("failed to set _dnslink (%s)", res.Status)
	}

	return nil;
}

func Update(ctx iris.Context) {
	body, ok := ctx.Values().Get("JSONBody").(iris.Map)
	if !ok {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("Bad request body"))
	}

	if ref, ok := body["ref"]; !ok {
		log.Println("Request Recieved: no reference key")
		_, _ = ctx.JSON(iris.Map{"message": "successful request"})
	} else if ref.(string) != "refs/heads/main" {
		log.Printf("Request Recieved: branch is not main (%v)\n", ref.(string))
		_, _ = ctx.JSON(iris.Map{"message": "successful request"})
	} else {
		cfg, ok := ctx.Values().Get("config").(*api39.Config)
		if !ok {
			log.Println("failed to retrieve configuration from context")
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"message": "failed to retrieve configuration from context",
			})
			return
		}

		if err := api39.UpdateGitRepo(cfg.Site.Repo, cfg.Site.Path); err != nil {
			message := "failed to update repository"
			log.Printf("%s: %v\n", message, err)
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"message": message,
				"error":   err,
			})
			return
		}

		if err := api39.RebuildWithMake(cfg.Site.Path); err != nil {
			message := "failed to rebuild site"
			log.Printf("%s: %v\n", message, err)
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"message": message,
				"error":   err,
			})
			return
		}

		path := filepath.Join(cfg.Site.Path, "public")
		hash, err := api39.IpfsAddDir(cfg.Ipfs.Url, path)
		if err != nil {
			message := "failed to add site to IPFS"
			log.Printf("%s: %v\n", message, err)
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"message": message,
				"error":   err,
			})
			return
		}

		log.Printf("New IPFS Hash: %s\n", hash)

		err = api39.UpdateDNSLink(cfg, hash)
		if err != nil {
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"message": "failed to update DNS settings",
				"error": err,
			})
		}

		_, _ = ctx.JSON(iris.Map{"message": "successful", "hash": hash})
	}
}
