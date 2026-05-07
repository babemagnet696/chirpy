package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	buffer := make([]byte, 32)
	_, err := rand.Read(buffer)
	if err != nil {
		return ""
	}

	token := hex.EncodeToString(buffer)
	return token
}