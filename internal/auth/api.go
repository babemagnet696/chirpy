package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	key, found := strings.CutPrefix(authHeader, "ApiKey ")
	if !found {
		return "", fmt.Errorf("invalid header")
	}
	key = strings.TrimSpace(key)

	if key == "" {
		return "", fmt.Errorf("no token provided")
	}

	return key, nil
}