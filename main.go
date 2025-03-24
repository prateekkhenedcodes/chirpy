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
	const port = "8080"
	const filePathRoot = "."

	//new multiplexer
	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	apicfg := &apiConfig{}
	mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apicfg.CountHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.ResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirpHandler)

	//server configuration
	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("serving files from %s on port %s\n", filePathRoot, port)
	log.Fatal(s.ListenAndServe())
}