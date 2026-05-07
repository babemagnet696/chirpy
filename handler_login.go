package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/babemagnet696/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string        `json:"password"`
		Email            string        `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON", err)
		return
	}

	dbUser, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Email not found", err)
		return
	}


	expirationTime := 1 * time.Hour

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating session", err)
	}
	refresh_token := auth.MakeRefreshToken()	

	ok, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking login", err)
	}
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect login", nil)
		return
	}

	user := dbUsertoUser(dbUser)


	respondWithJSON(w, http.StatusOK, response{
		User: user,
		Token: token,
		RefreshToken: refresh_token,
	})
}