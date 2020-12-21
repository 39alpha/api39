package api39

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
)

func Pipeline(app *iris.Application, cfgfile string) error {

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
		body, err := ctx.GetBody()
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}

		var parsed iris.Map
		if err := json.Unmarshal(body, &parsed); err != nil {
			ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
				"error": "failed to parse request body",
			})
			return
		} else {
			ctx.Values().Set("JSONBody", parsed)
		}
	}
	ctx.Next()
}

func WithConfig(cfg *Config) iris.Handler {
	return func(ctx iris.Context) {
		ctx.Values().Set("config", cfg)
		ctx.Next()
	}
}

func VerifyGithubSignature(ctx iris.Context) {
	cfg, ok := ctx.Values().Get("config").(*Config)
	if !ok {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error": "no configuration loaded",
		})
		return
	}

	var headers struct {
		HubSignature string `header:"X-Hub-Signature-256,required"`
	}

	if err := ctx.ReadHeaders(&headers); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	signature := headers.HubSignature[7:len(headers.HubSignature)]

	body, err := ctx.GetBody()
	if err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"error": "cannot read request body",
		})
		return
	}

	if !IsValidMAC(body, []byte(signature), []byte(cfg.Apikey)) {
		ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{
			"error": "invalid request signature",
		})
	}

	ctx.Next()
}
