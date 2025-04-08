package cmd

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func CreateBin(s *state.State, flags map[string]struct{}, args []string) error {
	p, v := utils.ValidateFlags(flags, "p"), utils.ValidateFlags(flags, "v")

	pathSlice := utils.ParseInputPath(args[0])

	last := pathSlice[len(pathSlice)-1]

	if len(pathSlice) > 1 {

		parentID := uuid.NullUUID{Valid: false}
		for i, e := range pathSlice {
			lastEle := i == len(pathSlice)-1
			bin, err := s.DBQueries.GetBin(context.Background(), database.GetBinParams{
				Name:     e,
				ParentID: parentID,
			})
			if err != nil {
				if !lastEle && !p {
					return fmt.Errorf("cannot create bin '%s': No such parent bin", e)
				}

				newBin, err := s.DBQueries.CreateBin(context.Background(), database.CreateBinParams{
					Name:     e,
					ParentID: parentID,
				})
				if err != nil {
					return fmt.Errorf("issue creating bin '%s': %v", e, err)
				}

				parentID = uuid.NullUUID{Valid: true, UUID: newBin.ID}

				if v {
					fmt.Printf("created bin '%s'\n", newBin.Name)
				}

			} else {
				parentID = uuid.NullUUID{Valid: true, UUID: bin.ID}

				if v {
					fmt.Printf("cannot create bin '%s': bin exists\n", bin.Name)
				}

			}
		}

		return nil
	}

	bin, err := s.DBQueries.CreateBin(context.Background(), database.CreateBinParams{
		Name:     last,
		ParentID: uuid.NullUUID{Valid: false},
	})
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("created bin '%s'\n", bin.Name)
	}

	return nil
}

func GetBin(s *state.State, flags map[string]struct{}, args []string) error {
	srcSlice := utils.ParseInputPath(args[0])
	lastElem, err := utils.GetLastElement(s, srcSlice)
	if err != nil {
		return err
	}

	bins, err := s.DBQueries.GetBinsByParent(context.Background(), uuid.NullUUID{
		Valid: true,
		UUID:  lastElem.ID.UUID,
	})
	if err != nil {
		return err
	}

	for _, e := range bins {
		fmt.Println(e.Name)
	}

	return nil
}
