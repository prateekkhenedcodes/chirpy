package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/prateekkhenedcodes/chirpy/internal/auth"
)

func (cfg *apiConfig) WebhookHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "request body could not decode", err)
		return
	}

	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "could not get the apikey", err)
		return
	}
	if apikey != cfg.polkaKey {
		respondWithError(w, 401, "apikeys have not matched", fmt.Errorf("apikeys have not matched"))
	}
	if params.Event != "user.upgraded" {
		respondWithError(w, 204, "not a valid event", err)
		return
	}
	UserUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse the userid string to UUID", err)
		return
	}

	_, err = cfg.dbQ.UpgradeUserRed(r.Context(), UserUUID)
	if err != nil {
		respondWithError(w, 404, "user not found or error during updating the chirp_red", err)
		return
	}

	respondWithJSON(w, 204, nil)

}
