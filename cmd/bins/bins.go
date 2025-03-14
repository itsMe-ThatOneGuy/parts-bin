package bins

import (
	"context"
	"errors"
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
					msg := fmt.Sprintf("mkbin: cannot create bin '%s': no such parent bin", e)
					return errors.New(msg)
				}

				newBin, err := s.DBQueries.CreateBin(context.Background(), database.CreateBinParams{
					Name:     e,
					ParentID: parentID,
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

	bin, err := s.DBQueries.CreateBin(context.Background(), database.CreateBinParams{
		Name:     last,
		ParentID: uuid.NullUUID{Valid: false},
	})
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("bin '%s' created\n", bin.Name)
	}

	return nil
}

func GetBin(s *state.State, flags map[string]struct{}, args []string) error {
	sourceSlice := utils.ParseInputPath(args[0])

	sourceParentID := uuid.NullUUID{Valid: false}
	lastBinInSource := database.Bin{}
	for _, e := range sourceSlice {
		bin, err := s.DBQueries.GetBin(context.Background(), database.GetBinParams{
			Name:     e,
			ParentID: sourceParentID,
		})
		if err != nil {
			return err
		}

		sourceParentID = uuid.NullUUID{Valid: true, UUID: bin.ID}

		lastBinInSource = bin
	}

	bins, err := s.DBQueries.GetBinsByParent(context.Background(), uuid.NullUUID{
		Valid: true,
		UUID:  lastBinInSource.ID,
	})
	if err != nil {
		return err
	}

	for _, e := range bins {
		fmt.Println(e.Name)
	}

	return nil
}
