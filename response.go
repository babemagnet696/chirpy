package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVal struct {
		Body  string `json:"error"`
	}

	var errorVal returnVal
	errorVal.Body = msg
	resp, err := json.Marshal(&errorVal)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}