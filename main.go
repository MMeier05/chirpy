package main

import (
	"net/http"
	"log"
	"sync/atomic"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
const port = "8080"
const filepathRoot = "."

func main() {
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetRequestsHandler)
	mux.HandleFunc("POST /api/validate_chirp", checkValidHandler)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	// Start the server 
	log.Fatal(srv.ListenAndServe())
}
