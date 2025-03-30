package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prateekkhenedcodes/chirpy/internal/auth"
)

func (cfg *apiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	type retToken struct {
		Token string `json:"token"`
	}
	authToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, "Authoraization header error", err)
		return
	}

	refreshToken, err := cfg.dbQ.GetRefreshToken(r.Context(), authToken)
	if err != nil {
		respondWithError(w, 401, "could not get the refresh token data table", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "token has expired", fmt.Errorf("token has expired"))
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "refresh token is revoked", fmt.Errorf("refresh token is revoked"))
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not make a new JWT", err)
		return
	}

	respondWithJSON(w, 200, retToken{
		Token: newToken,
	})
}
