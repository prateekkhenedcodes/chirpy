package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) ChirpGetHandler(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 404, "Not a vaid ID", err)
		return
	}

	data, err := cfg.dbQ.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, "chirpId does not exists", err)
		return
	}

	respondWithJSON(w, 200, ChirpReturnVals{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Body:      data.Body,
		UserId:    data.UserID,
	})
}
