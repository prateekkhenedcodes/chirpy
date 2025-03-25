package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
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

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_body: clean,
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
