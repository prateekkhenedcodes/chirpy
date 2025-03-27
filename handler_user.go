package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) CreatUserHandler(w http.ResponseWriter, r *http.Request) {

	type Parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	userData, err := cfg.dbQ.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudldn't create the user table in database", err)
		return
	}

	respondWithJSON(w, 201, User{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
	})

}
