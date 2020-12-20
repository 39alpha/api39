package config

import (
	"context"
	"github.com/39alpha/api39/auth"
	"net/http"
)

type WithConfig struct {
	cfg     *Config
	handler http.Handler
}

func (wc *WithConfig) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.WithValue(req.Context(), "config", wc.cfg)
	reqWithCfg := req.WithContext(ctx)
	wc.handler.ServeHTTP(w, reqWithCfg)
}

func NewWithConfig(filename string, handler http.Handler) (*WithConfig, error) {
	cfg, err := ReadConfig(filename)
	if err != nil {
		return nil, err
	}

	authed := auth.NewEnsureAuth(cfg.Apikey, handler)
	return &WithConfig{cfg, authed}, nil
}
