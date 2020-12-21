package site

import (
	"github.com/39alpha/api39/api39"
	"github.com/kataras/iris/v12"
)

func Update(ctx iris.Context) {
	if cfg, ok := ctx.Values().Get("config").(*api39.Config); !ok {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "no configuration loaded",
		})
	} else {
		ctx.JSON(iris.Map{"apikey": cfg.Apikey})
	}
}
