package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prateekkhenedcodes/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQ            *database.Queries
	platform       string
}

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println(err)
	}
	apicfg := &apiConfig{}
	apicfg.platform = os.Getenv("PLATFORM")
	if apicfg.platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	apicfg.dbQ = database.New(db)

	const port = "8080"
	const filePathRoot = "."

	//new multiplexer
	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))

	mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apicfg.CountHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.ResetHandler)
	mux.HandleFunc("POST /api/users", apicfg.CreatUserHandler)
	mux.HandleFunc("POST /api/chirps", apicfg.ChirpHandler)
	mux.HandleFunc("GET /api/chirps", apicfg.ChirpsGetHandler)


	//server configuration
	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("serving files from %s on port %s\n", filePathRoot, port)
	log.Fatal(s.ListenAndServe())
}
