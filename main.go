package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/SicklesScript/portfoliotrackerV2/internal/database"
)

type apiConfig struct {
	dbQueries *database.Queries
	jwtSecret string
	fmpApiKey string
	platform  string
}

func main() {
	const filePathRoot = "./static"
	const port = "8080"
	mux := http.NewServeMux()

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	queries := database.New(db)

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		fmt.Printf("PLATFORM must be set\n")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Printf("JWT_SECRET must be set\n")
	}

	fmpApiKey := os.Getenv("FMP_API_KEY")
	if fmpApiKey == "" {
		fmt.Printf("FMP_API_KEY must be set\n")
	}

	apiCfg := apiConfig{
		dbQueries: queries,
		jwtSecret: jwtSecret,
		fmpApiKey: fmpApiKey,
		platform:  platform,
	}

	fs := http.FileServer((http.Dir(filePathRoot)))
	mux.Handle("/", fs)
	mux.HandleFunc("/login", apiCfg.handlerLoginView)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginAPI)

	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Serving on port: %s\n", port)
	err = s.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
