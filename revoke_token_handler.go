package main

import (
	"net/http"

	"github.com/prateekkhenedcodes/chirpy/internal/auth"
)

func (cfg *apiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnot get the token from auth header", err)
		return
	}

	_, err = cfg.dbQ.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke refresh token", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
