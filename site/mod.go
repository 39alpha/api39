package site

import (
	"fmt"
	"github.com/39alpha/api39/config"
	"net/http"
)

func Update(w http.ResponseWriter, req *http.Request) {
	if cfg, ok := req.Context().Value("config").(*config.Config); !ok {
		fmt.Fprintf(w, "No configuration found\n")
	} else {
		fmt.Fprintf(w, "Apikey: %q\n", cfg.Apikey)
	}
}
