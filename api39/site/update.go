package site

import (
	"fmt"
	"github.com/39alpha/api39/api39"
	"github.com/kataras/iris/v12"
	"log"
	"path/filepath"
)

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
		_, _ = ctx.JSON(iris.Map{"message": "successful", "hash": hash})
	}
}
