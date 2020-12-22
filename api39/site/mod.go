package site

import (
	"fmt"
	"github.com/39alpha/api39/api39"
	"github.com/kataras/iris/v12"
	"log"
	"os"
)

func Update(ctx iris.Context) {
	body, ok := ctx.Values().Get("JSONBody").(iris.Map)
	if !ok {
		ctx.StopWithError(iris.StatusBadRequest, fmt.Errorf("Bad request body"))
	}

	if ref, ok := body["ref"]; !ok {
		log.Println("Request Recieved: no reference key")
	} else if ref.(string) != "refs/heads/main" {
		log.Printf("Request Recieved: branch is not main (%v)\n", ref.(string))
	} else {
		cfg, ok := ctx.Values().Get("config").(*api39.Config)
		if !ok {
			log.Println("Failed to retrieve configuration from context")
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"error": "an error occurred while loading site configuration",
			})
			return
		}

		if err := api39.UpdateGitRepo(cfg.Site.Repo, cfg.Site.Path); err != nil {
			log.Printf("Failed to update repository: %v\n", err)
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"message": "failed to update repository",
				"error":   err,
			})
			return
		}

		if err := api39.RebuildWithHugo(cfg.Site.Hugo, cfg.Site.Path); err != nil {
			log.Printf("Failed to rebuild site: %v\n", err)
			ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
				"message": "failed to rebuild site",
				"error":   err,
			})
			return
		}
		cwd, _ := os.Getwd()
		log.Printf("Working Directory: %v\n", cwd)
	}

	ctx.JSON(iris.Map{"message": "successful request"})
}
