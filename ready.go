package main

import "net/http"

func greetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func (cfg *apiConfig) getAppHandler(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	stripped := http.StripPrefix("/app", fs)
	return cfg.middlewareMetricsInc(stripped)
}