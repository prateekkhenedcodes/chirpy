package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/chirpy/internal/database"
)
type ChirpReturnVals struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}
func (cfg *apiConfig) ChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string `json:"body"`
		User_id uuid.UUID
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	clean := profaneClean(params.Body)

	dbData, err := cfg.dbQ.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   clean,
		UserID: params.User_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp post", err)
		return
	}

	respondWithJSON(w, 201, ChirpReturnVals{
		ID:        dbData.ID,
		CreatedAt: dbData.CreatedAt,
		UpdatedAt: dbData.UpdatedAt,
		Body:      dbData.Body,
		UserId:    dbData.UserID,
	})
}

func profaneClean(s string) string {
	pwords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	sslice := strings.Split(s, " ")

	for i, word := range sslice {
		for _, pword := range pwords {
			if strings.ToLower(word) == pword {
				sslice[i] = "****"
			}
		}
	}
	return strings.Join(sslice, " ")
}
