package bins

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func CreateBin(s *state.State, flags map[string]struct{}, args []string) error {
	p, v := utils.ValidateFlags(flags, "p"), utils.ValidateFlags(flags, "v")

	last, _, pathSlice := utils.ParseInputPath(args[0])
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

func DeleteBin(s *state.State, flags map[string]struct{}, args []string) error {
	r, v := utils.ValidateFlags(flags, "r"), utils.ValidateFlags(flags, "v")

	last, _, pathSlice := utils.ParseInputPath(args[0])

	if r {
		pathCache := make(map[string]struct {
			Name   string
			ID     uuid.UUID
			Parent uuid.NullUUID
		})

		parentID := uuid.NullUUID{Valid: false}
		for _, e := range pathSlice {
			bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
				Name:      e,
				ParentBin: parentID,
			})
			if err != nil {
				return err
			}

			pathCache[e] = struct {
				Name   string
				ID     uuid.UUID
				Parent uuid.NullUUID
			}{
				Name:   e,
				ID:     bin.ID,
				Parent: parentID,
			}

			parentID = uuid.NullUUID{Valid: true, UUID: bin.ID}
		}

		slices.Reverse(pathSlice)
		for i, e := range pathSlice {
			bin := pathCache[e]

			bins, err := s.DBQueries.GetBinsByParent(context.TODO(), uuid.NullUUID{
				Valid: true,
				UUID:  bin.ID,
			})
			if err != nil {
				return err
			}

			if len(bins) > 1 {
				for _, e := range bins {
					if v {
						fmt.Printf("deleting '%s'\n", e.Name)
					}

					_, err := s.DBQueries.DeleteBin(context.TODO(), database.DeleteBinParams{
						Name:      e.Name,
						ParentBin: uuid.NullUUID{Valid: true, UUID: bin.ID},
					})
					if err != nil {
						return err
					}
				}
			}

			if i != len(pathSlice)-1 {
				if v {
					fmt.Printf("deleting '%s'\n", bin.Name)
				}

				_, err = s.DBQueries.DeleteBin(context.TODO(), database.DeleteBinParams{
					Name:      bin.Name,
					ParentBin: bin.Parent,
				})
				if err != nil {
					return err
				}
			}

		}

		return nil
	}

	if v {
		fmt.Printf("deleting '%s'\n", last)
	}

	_, err := s.DBQueries.DeleteBin(context.TODO(), database.DeleteBinParams{
		Name:      last,
		ParentBin: uuid.NullUUID{Valid: false},
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateBin(s *state.State, flags map[string]struct{}, args []string) error {
	v := utils.ValidateFlags(flags, "v")

	_, _, sourceSlice := utils.ParseInputPath(args[0])
	_, _, destinationSlice := utils.ParseInputPath(args[1])

	sourceParentID := uuid.NullUUID{Valid: false}
	lastBinInSource := database.Bin{}
	for _, e := range sourceSlice {
		bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
			Name:      e,
			ParentBin: sourceParentID,
		})
		if err != nil {
			return err
		}

		sourceParentID = uuid.NullUUID{Valid: true, UUID: bin.ID}

		lastBinInSource = bin
	}

	destinationParentID := uuid.NullUUID{Valid: false}
	for i, e := range destinationSlice {
		bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
			Name:      e,
			ParentBin: destinationParentID,
		})
		if err != nil {
			if i != len(destinationSlice)-1 {
				return err
			}

			bin, err := s.DBQueries.UpdateBinName(context.Background(), database.UpdateBinNameParams{
				Name:      lastBinInSource.Name,
				ParentBin: destinationParentID,
				Name_2:    e,
			})
			if err != nil {
				return nil
			}

			lastBinInSource = bin

			break
		}

		destinationParentID = uuid.NullUUID{Valid: true, UUID: bin.ID}

	}

	err := s.DBQueries.UpdateBinParent(context.Background(), database.UpdateBinParentParams{
		Name:        lastBinInSource.Name,
		ParentBin:   lastBinInSource.ParentBin,
		ParentBin_2: destinationParentID,
	})
	if err != nil {
		return err
	}

	if v {
		msg := fmt.Sprintf("renamed '%v' -> '%v'", args[0], args[1])
		if len(args[1]) < len(args[0]) {
			msg = fmt.Sprintf("renamed '%v' -> '%v/%v'", args[0], args[1], lastBinInSource.Name)
		}

		fmt.Println(msg)

	}

	return nil
}

func GetBin(s *state.State, flags map[string]struct{}, args []string) error {
	_, _, sourceSlice := utils.ParseInputPath(args[0])

	sourceParentID := uuid.NullUUID{Valid: false}
	lastBinInSource := database.Bin{}
	for _, e := range sourceSlice {
		bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
			Name:      e,
			ParentBin: sourceParentID,
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
