package main

import (
	"net/http"
	"strings"
	"fmt"
)

func validate(w http.ResponseWriter, body string) (string, error) {

	if len(body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return body, fmt.Errorf("Chirp is too long")
	}

	return profanityCheck(body), nil

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