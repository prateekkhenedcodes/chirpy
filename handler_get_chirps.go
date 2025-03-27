package main

import (
	"net/http"
)

func (cfg *apiConfig) ChirpsGetHandler(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.dbQ.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get all the chirps", err)
		return
	}

	chirpReturnVals := make([]ChirpReturnVals, 0, len(chirps))

	for _, chirp := range chirps {
		chirpReturnVals = append(chirpReturnVals, ChirpReturnVals{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	respondWithJSON(w, 200, chirpReturnVals)

}
