package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"sort"
	"github.com/babemagnet696/chirpy/internal/auth"
	"github.com/babemagnet696/chirpy/internal/database"
	"github.com/google/uuid"
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
		Body     string  `json:"body"`
		AuthorID string  `json:"author_id"`
	}

	userID, err := cfg.getIdFromRequest(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed JSON", err)
		return
	}


	cleanRes, err := validateChirpBody(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	params.Body = cleanRes

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:     params.Body,
		UserID:   userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	chirp := dbChirptoChirp(dbChirp)

	respondWithJSON(w, http.StatusCreated, chirp)
}


func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	sUserID := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")
	switch sortBy {
		case "asc":
			break
		case "desc":
			break
		default:
			sortBy = "asc"
	}

	// Check for optional userID in URL
	// Return chirps by user_id instead
	if sUserID != "" {
		userID, err := uuid.Parse(sUserID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing userID", err)
			return
		}

		resp, err := cfg.db.GetChirpsByUser(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "user not found", err)
			return
		}

		chirps := dbChirpSlicetoChirpSlice(resp)
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	// return all chirps
	resp, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps", nil)
		return
	}
	

	chirps := dbChirpSlicetoChirpSlice(resp)
	if sortBy == "asc" {
		sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)})
	}
	if sortBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.After(chirps[j].CreatedAt)})
	}
	

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


func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
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
	


	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}
	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "forbidden", err)
		return
	}
	

	if err = cfg.db.DeleteChirp(r.Context(), chirpID) ; err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
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