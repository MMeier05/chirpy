package main

import (
	"time"
	"net/http"
	"encoding/json"
	"database/sql"
	"github.com/MMeier05/chirpy/internal/auth"
	"github.com/MMeier05/chirpy/internal/database"
)
type Parameters struct {
	Password string `json:"password"`
	Email string `json:"email"`
	ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}

	//Call to Database to store user  
	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPW,
			Valid: true,
		},
	})
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

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	isSame, err := auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if !isSame {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	reqExpirationTime := time.Duration(1 * time.Hour) 
	if params.ExpiresInSeconds != 0 && params.ExpiresInSeconds < 3600 {
		reqExpirationTime = time.Duration(params.ExpiresInSeconds * time.Second)
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, reqExpirationTime)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	resp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}
	err = respondWithJSON(w, 200, resp)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
}
