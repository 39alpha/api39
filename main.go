package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/39alpha/api39/site"
	"log"
	"net/http"
)

const apikeylen = 64

var (
	port   = 3964
	genkey = false
)

func init() {
	flag.IntVar(&port, "port", port, "port on which the server will listen")
	flag.BoolVar(&genkey, "genkey", genkey, "generate and print a random API key to STDOUT and exit")
}

func generateApiKey() (string, error) {
	chars := [64]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
		'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
		'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '+',
		'/'}

	key := make([]byte, apikeylen)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	for i, x := range key {
		key[i] = chars[x%64]
	}
	return string(key), nil
}

func main() {
	flag.Parse()

	if genkey {
		apikey, err := generateApiKey()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(apikey)
	} else {
		addr := fmt.Sprintf(":%d", port)
		http.HandleFunc("/api/v0/site/update", site.Update)
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}
