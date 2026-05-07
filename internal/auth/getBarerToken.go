package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	
	token, found := strings.CutPrefix(authHeader, "Bearer ")
	if !found {
		return "", fmt.Errorf("invalid header")
	}
	token = strings.TrimSpace(token)

	if token == "" {
		return "", fmt.Errorf("no token provided")
	}

	return token, nil
}