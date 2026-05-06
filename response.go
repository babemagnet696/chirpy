package main

import (
	"net/http"
	"log"
	"encoding/json"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type errorResponse struct {
		Error  string `json:"error"`
	}

	if err != nil {
        log.Println(err)
    }

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	var errorVal errorResponse
	errorVal.Error = msg
	respondWithJSON(w, code, errorVal)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
    	w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(resp)
}