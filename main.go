package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/config"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		return
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	godotenv.Load(cfg.EVNPath)
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = cfg.DBUrl
		if dbURL == "" || dbURL == config.DefaultDBurl {
			log.Fatal("db_url not set in either .env or ~/.partsbinconfig.json")
		}
	}

	dbCon, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %s", err)
	}
	defer dbCon.Close()

}

func server() {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + "8080",
	}

	log.Println("Serving on port 8080")
	log.Fatal(server.ListenAndServe())
}
