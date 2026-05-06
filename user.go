package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"time"
	"github.com/babemagnet696/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	user := dbUsertoUser(dbUser)
	
	respondWithJSON(w, http.StatusCreated, user)
}

func dbUsertoUser(du database.User) User {
	return User{
		ID:        du.ID,
		CreatedAt: du.CreatedAt,
		UpdatedAt: du.UpdatedAt,
		Email:     du.Email,
	}
}