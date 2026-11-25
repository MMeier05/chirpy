package main

import (
	"database/sql"
	"errors"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/MMeier05/chirpy/internal/database"
	"github.com/google/uuid"
)

	type returnVals struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

func (cfg *apiConfig) createPostHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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

	cleanedString, err := replaceProfanity(params.Body)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	post, err := cfg.db.CreatePost(req.Context(), database.CreatePostParams{
		Body: cleanedString,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	vals := returnVals{
		ID: post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body: post.Body,
		UserID: post.UserID,
	}
	err = respondWithJSON(w, 201, vals)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
}

func (cfg *apiConfig) getAllPostsHandler(w http.ResponseWriter, req *http.Request) {
	posts, err := cfg.db.GetAllPosts(req.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	postsJSON := make([]returnVals, len(posts))
	for i, post := range posts {
		postsJSON[i] = returnVals{
			ID: post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Body: post.Body,
			UserID: post.UserID,
		}
	}
	err = respondWithJSON(w, http.StatusOK, postsJSON)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
}
func (cfg *apiConfig) getPostHandler(w http.ResponseWriter, req *http.Request) {
	userIdString := req.PathValue("chirpID")
	userID, err := uuid.Parse(userIdString)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	post, err := cfg.db.GetPost(req.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "User not found")
			return 
		} else {
		respondWithError(w, 500, "Something went wrong2")
		return
		}
	}
	postJSON := returnVals{
			ID: post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Body: post.Body,
			UserID: post.UserID,
		}
	err = respondWithJSON(w, http.StatusOK, postJSON)
	if err != nil {
		respondWithError(w, 500, "Something went wrong3")
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

