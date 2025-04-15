package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func Repl(s *state.State) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("parts-bin > ")

		scanner.Scan()
		input := parseInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		cmdName := input[0]
		if cmdName == "exit" {
			break
		}

		hasFlags := false
		flags := make(map[string]string)
		if len(input) > 1 {
			if strings.HasPrefix(input[1], "-") {
				flags = utils.ParseFlags(input, &hasFlags)
			}
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
			continue
		} else {
			fmt.Printf("command not found: %s\n", cmdName)
			continue
		}

	}
}

func parseInput(input string) []string {
	return strings.Fields(input)
}
