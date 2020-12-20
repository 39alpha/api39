package site

import (
    "fmt"
    "net/http"
)

func Update(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Updating the site still")
}
