package main

import (
	"flag"
	"fmt"
	"github.com/39alpha/api39/config"
	"github.com/39alpha/api39/site"
	"log"
	"net/http"
	"os"
)

const apikeylen = 64

var (
	port       = 3964
	genconf    = false
	configpath = ""
)

func init() {
	flag.IntVar(&port, "port", port, "port on which the server will listen")
	flag.BoolVar(&genconf, "genconf", genconf, "generate and print a configuration file to STDOUT and exit")
	flag.StringVar(&configpath, "config", configpath, "path to configuration file (required)")
}

func main() {
	flag.Parse()

	if genconf {
		err := config.GenerateConfig(apikeylen)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if configpath == "" {
			fmt.Fprintf(os.Stderr, "Error: -config flag is required\n\n")
			flag.Usage()
			os.Exit(1)
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/api/v0/site/update", site.Update)

		api, err := config.NewWithConfig(configpath, mux)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		addr := fmt.Sprintf(":%d", port)
		log.Fatal(http.ListenAndServe(addr, api))
	}
}
