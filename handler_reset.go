package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not authorized", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error resetting server", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0\n"))
	w.Write([]byte("All data deleted\n"))

}