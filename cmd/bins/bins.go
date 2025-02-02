package bins

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

func validateInputType(s string) (string, uuid.UUID) {
	inputType := "UUID"
	uuidTest, err := uuid.Parse(s)
	if err != nil {
		inputType = "string"
		return inputType, uuid.Nil
	}

	return inputType, uuidTest
}

func validateInputPath(s string) (last string, parent string, pathSlice []string) {
	splitSlice := strings.Split(s, "/")
	if splitSlice[0] == "" {
		splitSlice = splitSlice[1:]
	}

	lastIndex := len(splitSlice) - 1
	_last := splitSlice[lastIndex]

	if len(splitSlice) > 1 {
		parentIndex := lastIndex - 1
		_parent := splitSlice[parentIndex]
		return _last, _parent, splitSlice
	}

	return _last, "", splitSlice
}

func CreateBin(s *state.State, args []string) (string, error) {
	last, _, pathSlice := validateInputPath(args[0])
	if len(pathSlice) > 1 {
		for _, v := range pathSlice {
			_, err := s.DBQueries.GetBinByName(context.TODO(), v)
			if err != nil {
				msg := fmt.Sprintf("mkbin: cannot create bin '%s': no such parent bin", v)
				return "", errors.New(msg)
			}
		}
	}

	_, err := s.DBQueries.CreateBin(context.TODO(), database.CreateBinParams{
		Name:      last,
		ParentBin: uuid.NullUUID{Valid: false},
	})
	if err != nil {
		msg := fmt.Sprintf("issue creating '%s' bin", last)
		return "", errors.New(msg)
	}

	msg := fmt.Sprintf("bin '%s' created", last)
	return msg, nil
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
