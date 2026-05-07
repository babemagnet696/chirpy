package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/babemagnet696/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	var apiCfg apiConfig

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	apiCfg.platform = os.Getenv("PLATFORM")


	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}


	databaseQueries := database.New(db)
	apiCfg.db = databaseQueries


	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.getAppHandler(filepathRoot))

	mux.HandleFunc("GET  /api/healthz",          greetHandler)
	mux.HandleFunc("GET  /admin/metrics",        apiCfg.handlerMetrics)
	mux.HandleFunc("GET  /api/chirps",           apiCfg.handlerGetChirps)
	mux.HandleFunc("GET  /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users",   apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login",   apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/chirps",  apiCfg.handlerChirp)
	

	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}


