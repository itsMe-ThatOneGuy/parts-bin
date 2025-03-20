package cmd

import (
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

type Command struct {
	Name        string
	Description string
	Callback    func(*state.State, map[string]struct{}, []string) error
}

var commands = []Command{
	{
		Name:        "bin",
		Description: "Create a bin in provided path",
		Callback:    CreateBin,
	},
	{
		Name:        "part",
		Description: "Create a part in provided path",
		Callback:    CreatePart,
	},
	{
		Name:        "ls",
		Description: "List parts and bins in provided path",
		Callback:    Ls,
	},
	{
		Name:        "mv",
		Description: "move a part/bin from provided source path to provided destination path",
		Callback:    Mv,
	},
	{
		Name:        "rm",
		Description: "remove part/bin in provided path",
		Callback:    Rm,
	},
}
