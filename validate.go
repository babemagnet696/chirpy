package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		Cleaned_body string `json:"clean_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	params.Cleaned_body = profanityCheck(params.Body)

	type successResult struct {
		Valid bool `json:"valid"`
	}
	respondWithJSON(w, http.StatusOK, successResult{Valid: true})

}

func profanityCheck(text string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	wordList := strings.Split(text, " ")

	for i, word := range wordList {
		cleanWord := strings.ToLower(word)

		if _, ok := badWords[cleanWord]; ok {
			wordList[i] = "****"
		}

	}

	cleanedText := strings.Join(wordList, " ")
	return cleanedText
}