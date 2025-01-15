package bins

import (
	"context"
	"log"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

func CreateBin(s *state.State, args []string) (database.Bin, error) {
	bin, err := s.DBQueries.CreateBin(context.Background(), args[0])
	if err != nil {
		log.Fatal("Error creating bin")
	}

	return bin, nil
}
