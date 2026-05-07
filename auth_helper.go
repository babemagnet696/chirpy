package main

import (
	"net/http"

	"github.com/babemagnet696/chirpy/internal/auth"
	"github.com/google/uuid"
)


func (cfg *apiConfig) getIdFromRequest(r http.Header) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(r)
	if err != nil {
		return uuid.Nil, err
	}
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}