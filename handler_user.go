package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"
	"github.com/babemagnet696/chirpy/internal/database"
	"github.com/babemagnet696/chirpy/internal/auth"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Email           string    `json:"email"`
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
	}
}