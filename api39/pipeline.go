package api39

import (
	"github.com/kataras/iris/v12"
)

func Pipeline(app *iris.Application, cfgfile string) error {
	withConfig, err := NewWithConfig(cfgfile)
	if err != nil {
		return err
	}

	app.UseGlobal(withConfig, ParseBody, EnsureAuth)

	return nil
}

func NewWithConfig(cfgfile string) (iris.Handler, error) {
	cfg, err := ReadConfig(cfgfile)
	if err != nil {
		return nil, err
	}

	return WithConfig(cfg), nil
}

func ParseBody(ctx iris.Context) {
	if ctx.Method() != "GET" {
		var body iris.Map
		if err := ctx.ReadJSON(&body); err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
				"error": "failed to parse request body",
			})
			return
		}
		ctx.Values().Set("JSONBody", body)
	}
	ctx.Next()
}

func WithConfig(cfg *Config) iris.Handler {
	return func(ctx iris.Context) {
		ctx.Values().Set("config", cfg)
		ctx.Next()
	}
}

func EnsureAuth(ctx iris.Context) {
	cfg, ok := ctx.Values().Get("config").(*Config)
	if !ok {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "no configuration loaded",
		})
		return
	}

	body := ctx.Values().Get("JSONBody").(iris.Map)
	if apikey, ok := body["apikey"]; !ok {
		ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{
			"error": "no api key provided",
		})
	} else if apikey.(string) != cfg.Apikey {
		ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{
			"error": "invalid api key",
		})
	} else {
		ctx.Next()
	}
}
