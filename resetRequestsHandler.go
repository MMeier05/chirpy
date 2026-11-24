package main

import (
	"net/http"
)

func (cfg *apiConfig) resetRequestsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
	w.WriteHeader(http.StatusOK)
}
