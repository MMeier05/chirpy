package main

import (
	"net/http"
	"encoding/json"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	type Parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	resp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	err = respondWithJSON(w, 201, resp)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
}

