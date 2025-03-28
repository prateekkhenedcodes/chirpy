package main

import (
	"encoding/json"
	"net/http"

	"github.com/prateekkhenedcodes/chirpy/internal/auth"
)

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
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

	respondWithJSON(w, 200, User{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
	})

}
