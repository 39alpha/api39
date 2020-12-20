package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
)

func GenerateApiKey(n int) (string, error) {
	chars := [64]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
		'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
		'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0',
		'1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '/'}

	key := make([]byte, n)

	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	for i, x := range key {
		key[i] = chars[x%64]
	}

	return string(key), nil
}

type Auth struct {
	Apikey string `json:token`
}

type EnsureAuth struct {
	Apikey  string
	handler http.Handler
}

func (ea *EnsureAuth) authenticate(req *http.Request) error {
	dec := json.NewDecoder(req.Body)

	var auth Auth
	if err := dec.Decode(&auth); err != nil {
		return fmt.Errorf("cannot parse request body: %v", err)
	}

	if auth.Apikey != ea.Apikey {
		return fmt.Errorf("invalid api key")
	}

	return nil
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := ea.authenticate(req); err != nil {
		fmt.Fprintf(w, "{ \"error\": %q }", err)
		return
	}
	ea.handler.ServeHTTP(w, req)
}

func NewEnsureAuth(apikey string, handler http.Handler) *EnsureAuth {
	return &EnsureAuth{apikey, handler}
}
