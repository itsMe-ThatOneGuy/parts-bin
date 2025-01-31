package bins

import (
	"context"
	"errors"

	"github.com/google/uuid"
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

func GetBin(s *state.State, args []string) (database.Bin, error) {
	argType := args[0]
	argInput := args[1]

	if argType == "-i" {
		argInput, err := uuid.Parse(argInput)
		if err != nil {
			return database.Bin{}, errors.New("Issue parsing UUID")
		}

		bin, err := s.DBQueries.GetBinByID(context.Background(), argInput)
		if err != nil {
			return database.Bin{}, errors.New("Issue getting bin using bin ID")
		}
		return bin, nil
	}

	if argType == "-n" {
		bin, err := s.DBQueries.GetBinByName(context.Background(), argInput)
		if err != nil {
			return database.Bin{}, errors.New("Issue getting bin using bin Name")
		}
		return bin, nil
	}

	return database.Bin{}, errors.New("Required argument flag not provided")
}

func DeleteBin(s *state.State, args []string) (database.Bin, error) {
	argType := args[0]
	argInput := args[1]

	if argType == "-i" {
		argInput, err := uuid.Parse(argInput)
		if err != nil {
			return database.Bin{}, errors.New("Issue parsing UUID")
		}

		bin, err := s.DBQueries.DeleteBinByID(context.Background(), argInput)
		if err != nil {
			return database.Bin{}, errors.New("Issue deleting bin using id")
		}

		return bin, nil
	}

	if argType == "-n" {
		bin, err := s.DBQueries.DeleteBinByName(context.Background(), argInput)
		if err != nil {
			return database.Bin{}, errors.New("Issue deleting bin using name")
		}

		return bin, nil
	}

	return database.Bin{}, errors.New("Required argument flag not provided")
}

func UpdateBin(s state.State, args []string) (database.Bin, error) {
	argType := args[0]
	argBin := args[1]
	argInputs := args[2:]

	if argType == "-i" {
		argBin, err := uuid.Parse(argBin)
		if err != nil {
			return database.Bin{}, errors.New("Issue parsing UUID")
		}

		if argInputs[0] == "-bn" {
			bin, err := s.DBQueries.UpdateBinNameByID(context.Background(),
				database.UpdateBinNameByIDParams{
					ID:   argBin,
					Name: argInputs[1],
				},
			)
			if err != nil {
				return database.Bin{}, errors.New("Issue renaming bin using id")
			}

			return bin, nil
		}

	}

	if argType == "-n" {
		if argInputs[0] == "-bn" {
			bin, err := s.DBQueries.UpdateBinNameByName(context.Background(),
				database.UpdateBinNameByNameParams{
					Name:   argBin,
					Name_2: argInputs[1],
				},
			)
			if err != nil {
				return database.Bin{}, errors.New("Issue renaming bin using name")
			}

			return bin, nil
		}

	}

	return database.Bin{}, nil
}
