package cmd

import (
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

type Command struct {
	Name        string
	Description string
	Callback    func(*state.State, map[string]string, []string) error
}

func Commands() map[string]Command {
	return map[string]Command{
		"help": {
			Name:        "help",
			Description: "Display the help message",
			Callback:    Help,
		},
		"mkbin": {
			Name:        "mkbin",
			Description: "Create a bin in provided path",
			Callback:    CreateBin,
		},
		"mkprt": {
			Name:        "mkprt",
			Description: "Create a part in provided path",
			Callback:    CreatePart,
		},
		"ls": {
			Name:        "ls",
			Description: "List parts and bins in provided path",
			Callback:    Ls,
		},
		"mv": {
			Name:        "mv",
			Description: "move a part/bin from provided source path to provided destination path",
			Callback:    Mv,
		},
		"rm": {
			Name:        "rm",
			Description: "remove part/bin in provided path",
			Callback:    Rm,
		},
	}

}
