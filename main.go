package main

import (
	"net/http"
	"log"
	"os"
	"time"
	"sync/atomic"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/MMeier05/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: database.New(db),
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllPostsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getPostHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetRequestsHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createPostHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginUserHandler)
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	// Start the server 
	log.Fatal(srv.ListenAndServe())
}
