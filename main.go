package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		return
	}

	state := &state.State{}

	err := state.InitConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	godotenv.Load(state.Config.EVNPath)

	err = state.InitDB()
	if err != nil {
		log.Fatalf("error connecting to db: %s", err)
	}
	defer state.CloseDB()


func CreateBin(s *state.State, args []string) (database.Bin, error) {
	bin, err := s.DBQueries.CreateBin(context.Background(), args[0])
	if err != nil {
		return database.Bin{}, err
	}

	return bin, nil
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
