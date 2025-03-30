package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/chirpy/internal/auth"
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

func (cfg *apiConfig) ChirpDeleteHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "could not get the access token", err)
		return
	}

	userIDJ, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "could not validate the access token ", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get the chirp id from the path", err)
		return
	}

	chirpData, err := cfg.dbQ.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}

	if chirpData.UserID != userIDJ {
		respondWithError(w, 403, "user is not the author of the chirp", err)
		return
	}

	err = cfg.dbQ.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete the chirp", err)
	}

	respondWithJSON(w, 204, nil)
}
