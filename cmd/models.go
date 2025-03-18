package cmd

import "github.com/itsMe-ThatOneGuy/parts-bin/internal/state"

type Command struct {
	Name        string
	Description string
	Callback    func(*state.State, map[string]struct{}, []string) error
}
