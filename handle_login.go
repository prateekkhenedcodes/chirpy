package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/chirpy/internal/auth"
	"github.com/prateekkhenedcodes/chirpy/internal/database"
)

type Login struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
	ReshreshToken string    `json:"refresh_token"`
	IsChirpRed    bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	token, err := auth.MakeJWT(userData.ID, cfg.jwtSecret, time.Duration(defaultExpTime)*time.Second)
	if err != nil {
		respondWithError(w, 401, "error during generating token string", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "refresh token was not generated", err)
		return
	}

	tokenData, err := cfg.dbQ.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userData.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create reshresh_token table ", err)
		return
	}

	dataB, err := cfg.dbQ.GetPass(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "email does not exists or error during getting data from the email ", err)
		return
	}

	respondWithJSON(w, 200, Login{
		ID:            userData.ID,
		CreatedAt:     userData.CreatedAt,
		UpdatedAt:     userData.UpdatedAt,
		Email:         userData.Email,
		Token:         token,
		ReshreshToken: tokenData.Token,
		IsChirpRed:    dataB.IsChirpyRed,
	})

}
