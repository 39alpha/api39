package api39

import (
	"crypto/rand"
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
