package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	var apiCfg apiConfig

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.getAppHandler(filepathRoot))

	mux.HandleFunc("GET  /api/healthz",        greetHandler)
	mux.HandleFunc("GET  /admin/metrics",      apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset",        apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal (srv.ListenAndServe())
}


