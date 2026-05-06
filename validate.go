package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed JSON", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}


	type successResult struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respondWithJSON(w, http.StatusOK, successResult{CleanedBody: profanityCheck(params.Body)})

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