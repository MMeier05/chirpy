package main

import (
	"net/http"
)

func (cfg *apiConfig) resetRequestsHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		respondWithError(w, 500, "Something unexpected happened")
		return
	}
	cfg.fileserverHits.Swap(0)
	w.WriteHeader(http.StatusOK)
}
