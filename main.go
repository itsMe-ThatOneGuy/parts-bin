package main

import (
	"log"
	"os"

	"github.com/itsMe-ThatOneGuy/parts-bin/cmd"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
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

	if len(os.Args) == 1 {
		cmd.Repl(state)
		return
	}

	cmd.RunCommand(state)

}
