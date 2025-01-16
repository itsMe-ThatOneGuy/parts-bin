package bins

import (
	"context"
	"errors"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

func CreateBin(s *state.State, args []string) (database.Bin, error) {
	bin, err := s.DBQueries.CreateBin(context.Background(), args[0])
	if err != nil {
		return database.Bin{}, errors.New("Issue creating new bin")
	}

	return bin, nil
}
