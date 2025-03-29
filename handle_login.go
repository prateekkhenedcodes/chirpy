package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/chirpy/internal/auth"
)

type Login struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password           string `json:"password"`
		Email              string `json:"email"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode the response body", err)
		return
	}

	userData, err := cfg.dbQ.GetPass(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get the password form the database", err)
		return
	}

	err = auth.CheckPasswordHash(userData.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	defaultExpTime := 3600

	if params.Expires_in_seconds > 0 {
		if params.Expires_in_seconds < defaultExpTime {
			defaultExpTime = params.Expires_in_seconds
		}
	}
	token, err := auth.MakeJWT(userData.ID, cfg.jwtSecret, time.Duration(defaultExpTime)*time.Second)
	if err != nil {
		respondWithError(w, 401, "error during generating token string", err)
	}

	respondWithJSON(w, 200, Login{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
		Token:     token,
	})

}
