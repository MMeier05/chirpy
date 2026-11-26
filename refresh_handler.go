package main

import (
	"time"
	"net/http"
	"github.com/MMeier05/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return 
	}
	params, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if time.Now().Compare(params.ExpiresAt) >= 0 || params.RevokedAt.Valid {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	newToken, err := auth.MakeJWT(params.UserID, cfg.jwtSecret, time.Duration(1 * time.Hour))
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	type returnVals struct {
		Token string `json:"token"`
	}
	resp := returnVals{
		Token : newToken,
	}
	err = respondWithJSON(w, 200, resp)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return 
	}
	err = cfg.db.RevokeToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	respondWithJSON(w, 204, nil)
}
