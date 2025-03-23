package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	//new multiplexer
	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	apicfg := &apiConfig{}
	mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /healthz", ReadinessHandler)
	mux.HandleFunc("GET /metrics", apicfg.CountHandler)
	mux.HandleFunc("POST /reset", apicfg.ResetHandler)

	//server configuration
	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("serving files from %s on port %s\n", filePathRoot, port)
	log.Fatal(s.ListenAndServe())
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) CountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Reset Successful"))
}
