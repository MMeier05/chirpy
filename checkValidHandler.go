package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func checkValidHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	cleanedString, err := replaceProfanity(params.Body)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	vals := returnVals{
		CleanedBody: cleanedString,
	}
	err = respondWithJSON(w, 200, vals)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
}

func replaceProfanity(body string) (string, error) {
	badWords := map[string]string{"kerfuffle": "", "sharbert": "", "fornax": ""}

	split := strings.Split(body, " ")
	loweredAndSplit := strings.Split(strings.ToLower(body), " ")

	for i, word := range loweredAndSplit {
		if _, ok := badWords[word]; ok {
			split[i] = "****"
			continue
		}
	}
	return strings.Join(split, " "), nil
}
