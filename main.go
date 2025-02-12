package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/itsMe-ThatOneGuy/parts-bin/cmd/bins"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		return
	}

	command := os.Args[1]
	input := os.Args[2:]

	hasFlags := false
	flags := make(map[string]struct{})
	if strings.HasPrefix(input[0], "-") {
		hasFlags = true
		_flags := strings.Split(input[0], "")[1:]
		for _, e := range _flags {
			flags[string(e)] = struct{}{}
		}
	}

	args := input
	if hasFlags {
		args = input[1:]
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

	switch command {
	case "createBin":
		bin, err := bins.CreateBin(state, args)
		if err != nil {
			log.Fatalf("Error creating bin: %v", err)
		}
		fmt.Println(bin.Name)
	case "getBin":
		bin, err := bins.GetBin(state, args)
		if err != nil {
			log.Fatalf("Error creating bin: %v", err)
		}
		fmt.Println(bin)
	case "deleteBin":
		bin, err := bins.GetBin(state, args)
		if err != nil {
			log.Fatalf("Error deleting bin: %v", err)
		}
		fmt.Println(bin)
	}

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
