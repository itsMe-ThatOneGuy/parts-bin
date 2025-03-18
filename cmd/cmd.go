package cmd

import (
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
)

var commands = []models.Command{
	{
		Name:        "bin",
		Description: "Create a bin in provided path",
		Callback:    CreateBin,
	},
}
