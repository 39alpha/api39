package main

import (
	"fmt"
	"github.com/39alpha/api39/site"
	"log"
	"net/http"
)

const port = 3964

var addr = fmt.Sprintf(":%d", port)

func main() {
	http.HandleFunc("/api/v0/site/update", site.Update)
	log.Fatal(http.ListenAndServe(addr, nil))
}
