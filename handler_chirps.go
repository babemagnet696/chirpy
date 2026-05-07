package main

import (
	"encoding/json"
	"net/http"
	"github.com/babemagnet696/chirpy/internal/database"
	"github.com/google/uuid"
	"time"
	"database/sql"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON", err)
		return
	}

	cleanRes, err := validate(w, params.Body)
	if err != nil {
		return // error is handled by validate already
	}

	params.Body = cleanRes

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:     params.Body,
		UserID:   params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	chirp := dbChirptoChirp(dbChirp)

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	resp, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps", nil)
		return
	}

	if len(resp) == 0 {
		emptyResponse := struct{
			Body string `json:"body"`
		}{
			Body: "There are no chirps to display",
		}
		respondWithJSON(w, http.StatusOK, emptyResponse)
		return
	}

	chirps := dbChirpSlicetoChirpSlice(resp)

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	resp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "No chirp found", err)
		return
	}
	chirp := dbChirptoChirp(resp)

	respondWithJSON(w, http.StatusOK, chirp)
}

func dbChirptoChirp(dc database.Chirp) Chirp {
	return Chirp{
		ID:        dc.ID,
		CreatedAt: dc.CreatedAt,
		UpdatedAt: dc.UpdatedAt,
		Body:      dc.Body,
		UserID:    dc.UserID,
	}
}

func dbChirpSlicetoChirpSlice(dc []database.Chirp) []Chirp {
	chirps := make([]Chirp, 0, len(dc))

	for _, c := range dc {
		chirps = append(chirps, dbChirptoChirp(c))
	}

	return chirps
}