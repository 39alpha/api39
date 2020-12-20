package site

import (
	"github.com/39alpha/api39/config"
	"github.com/39alpha/api39/respond"
	"net/http"
)

func Update(w http.ResponseWriter, req *http.Request) {
	if cfg, ok := req.Context().Value("config").(*config.Config); !ok {
		respond.ServerError(w)
	} else {
		respond.Respond(w, struct{ Apikey string }{cfg.Apikey})
	}
}
