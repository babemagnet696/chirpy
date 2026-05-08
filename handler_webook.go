package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/babemagnet696/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerRedChirpy(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	key, err := auth.GetApiKey(r.Header)
	if err != nil || key != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON", err)
		return
	}
	fmt.Printf("decoded params: %+v\n", params)
	fmt.Printf("event: %q\n", params.Event)
	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing id", err)
		return
	}
	if err = cfg.db.UpdateChirpyRed(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}