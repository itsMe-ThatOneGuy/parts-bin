package cmd

import (
	"fmt"
	"os"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func RunCommand(s *state.State) {
	input := os.Args[1:]
	cmdName := input[0]

	hasFlags := false
	flags := make(map[string]string)
	if len(input) > 1 {
		flags = utils.ParseFlags(input, &hasFlags)
	}

	args := input[1:]
	if hasFlags {
		args = input[2:]
	}

	command, ok := Commands()[cmdName]
	if ok {
		err := command.Callback(s, flags, args)
		if err != nil {
			fmt.Printf("%s: %v\n", cmdName, err)
		}
	} else {
		fmt.Println("Unknown Command")
	}
}
