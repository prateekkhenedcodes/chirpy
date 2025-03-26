package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden", errors.New("Forbidden"))
		return
	}

	cfg.dbQ.DeleteUsers(r.Context())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Reset Successful"))

}
