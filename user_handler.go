package main

import (
	"time"
	"net/http"
	"encoding/json"
	"database/sql"
	"github.com/google/uuid"
	"github.com/MMeier05/chirpy/internal/auth"
	"github.com/MMeier05/chirpy/internal/database"
)
type Parameters struct {
	Password 	string `json:"password"`
	Email 		string `json:"email"`
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
		Email: 					params.Email,
		HashedPassword: sql.NullString{
		String: 				hashedPW,
		Valid: 					true,
		},
	})
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}

	resp := User{
		ID: 				user.ID,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email: 			user.Email,
		IsChirpyRed: 	user.IsChirpyRed.Bool,
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

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(1 * time.Hour))
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	refreshToken := auth.MakeRefreshToken()
	cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: 			refreshToken,
		UserID: 		user.ID,
		ExpiresAt: 	time.Now().Add(time.Duration(24 * 60 * time.Hour)),
	})
	resp := User{
		ID: 					user.ID,
		CreatedAt: 		user.CreatedAt,
		UpdatedAt: 		user.UpdatedAt,
		Email: 				user.Email,
		Token: 				jwtToken,
		RefreshToken: refreshToken,
		IsChirpyRed: 	user.IsChirpyRed.Bool,
	}
	err = respondWithJSON(w, 200, resp)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
}

func (cfg *apiConfig) updateUserDataHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	err = cfg.db.UpdateUsersEmailPassword(req.Context(), database.UpdateUsersEmailPasswordParams{
		ID: userID,
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
	type returnVals struct {
		Email string `json:"email"`
	}
	vals := returnVals{
		Email: params.Email,
	}
	err = respondWithJSON(w, 200, vals)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
}

func (cfg *apiConfig) polkaWebhooksHandler(w http.ResponseWriter, req *http.Request) {
	type Parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	recievedKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if recievedKey != cfg.polkaKey {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}
	if params.Event == "user.upgraded" {
		err = cfg.db.UpgradeToChirpyRed(req.Context(), params.Data.UserID)
		if err != nil {
			respondWithError(w, 500, "Something unexpected happened")
			return
		}
	}
	err = respondWithJSON(w, 204, nil)
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}

}
