package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/babemagnet696/chirpy/internal/auth"
	"github.com/babemagnet696/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Email           string    `json:"email"`
	Token           string    `json:"token"`
	IsChirpyRed     bool      `json:"is_chirpy_red"`
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: hash,
		Email:          params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	user := dbUsertoUser(dbUser)
	
	respondWithJSON(w, http.StatusCreated, user)
}

func dbUsertoUser(du database.User) User {
	return User{
		ID:              du.ID,
		CreatedAt:       du.CreatedAt,
		UpdatedAt:       du.UpdatedAt,
		Email:           du.Email,
		IsChirpyRed:     du.IsChirpyRed.Bool,
	}
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	dbUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		HashedPassword: hash,
		Email:          params.Email,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error", err)
		return
	}

	user := dbUsertoUser(dbUser)
	respondWithJSON(w, http.StatusOK, user)
}