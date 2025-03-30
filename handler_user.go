package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/chirpy/internal/auth"
	"github.com/prateekkhenedcodes/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) CreatUserHandler(w http.ResponseWriter, r *http.Request) {

	type paraMeters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paraMeters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate hash for the password", err)
	}

	userData, err := cfg.dbQ.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudldn't create the user table in database with password", err)
		return
	}

	respondWithJSON(w, 201, User{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
	})

}

func (cfg *apiConfig) UpdateUsersHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "access token was not found in header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "token is invalid ", err)
		return
	}

	params := parameters{}
	decode := json.NewDecoder(r.Body)
	err = decode.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode the request body", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash the password", err)
		return
	}

	updatedData, err := cfg.dbQ.UpdateEmailPass(r.Context(), database.UpdateEmailPassParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update the email and password", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:        updatedData.ID,
		Email:     updatedData.Email,
		CreatedAt: updatedData.CreatedAt,
		UpdatedAt: updatedData.UpdatedAt,
	})

}
