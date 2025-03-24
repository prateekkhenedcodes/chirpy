package main

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Reset Successful"))
}

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}
	type errorRes struct {
		Error string `json:"error"`
	}
	type chirpRes struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resErr := errorRes{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(resErr)
		if err != nil {
			return
		}
		w.Write(dat)
		return
	}

	if len(req.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		resChirp := errorRes{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(resChirp)
		if err != nil {
			return
		}
		w.Write(dat)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		resChirp := chirpRes{
			Valid: true,
		}
		dat, err := json.Marshal(resChirp)
		if err != nil {
			return
		}
		w.Write(dat)
	}

}
