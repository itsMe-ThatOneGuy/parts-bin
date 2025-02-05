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

func validateFlags(flags map[string]struct{}, key string) bool {
	_, exists := flags[key]
	return exists
}

func CreateBin(s *state.State, flags map[string]struct{}, args []string) error {
	p, v := validateFlags(flags, "p"), validateFlags(flags, "v")

	last, _, pathSlice := validateInputPath(args[0])
	if len(pathSlice) > 1 {

		parentID := uuid.NullUUID{Valid: false}
		for _, e := range pathSlice {
			bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
				Name:      e,
				ParentBin: parentID,
			})
			if err != nil {
				if !p {
					msg := fmt.Sprintf("mkbin: cannot create bin '%s': no such parent bin", e)
					return errors.New(msg)
				}

				newBin, err := s.DBQueries.CreateBin(context.TODO(), database.CreateBinParams{
					Name:      e,
					ParentBin: parentID,
				})
				if err != nil {
					msg := fmt.Sprintf("issue creating '%s' bin: %v", e, err)
					return errors.New(msg)
				}

				parentID = uuid.NullUUID{Valid: true, UUID: newBin.ID}

				if v {
					fmt.Printf("bin '%s' created\n", newBin.Name)
				}

			} else {
				parentID = uuid.NullUUID{Valid: true, UUID: bin.ID}

				if v {
					fmt.Printf("bin '%s' already created\n", bin.Name)
				}

			}
		}

		return nil
	}

	bin, err := s.DBQueries.CreateBin(context.TODO(), database.CreateBinParams{
		Name:      last,
		ParentBin: uuid.NullUUID{Valid: false},
	})
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("bin '%s' created\n", bin.Name)
	}

	return nil
}


	if argType == "-n" {
		bin, err := s.DBQueries.GetBinByName(context.Background(), argInput)
		if err != nil {
			return database.Bin{}, errors.New("Issue getting bin using bin Name")
		}
		return bin, nil
	}
func DeleteBin(s *state.State, args []string) error {
	last, _, pathSlice := validateInputPath(args[0])
	if len(pathSlice) > 1 {

		parentID := uuid.NullUUID{Valid: false}
		for _, e := range pathSlice {
			_, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
				Name:      e,
				ParentBin: parentID,
			})
			if err != nil {
				return err
			}
		}

	}

	bin, err := s.DBQueries.DeleteBin(context.TODO(), database.DeleteBinParams{
		Name:      last,
		ParentBin: uuid.NullUUID{Valid: false},
	})
	if err != nil {
		return err
	}

	fmt.Printf("bin '%s' deleted\n", bin.Name)

	return nil
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
