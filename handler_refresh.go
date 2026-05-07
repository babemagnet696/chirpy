package main

import (
	"net/http"
	"time"
	"github.com/babemagnet696/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshTokenAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	newToken, err := auth.MakeJWT(user.ID, cfg.secret, 1 * time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating session token", err)
		return
	}
	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{newToken})
}