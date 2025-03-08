package bins

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
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
			bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
				Name:     e,
				ParentID: parentID,
			})
			if err != nil {
				if !lastEle && !p {
					msg := fmt.Sprintf("mkbin: cannot create bin '%s': no such parent bin", e)
					return errors.New(msg)
				}

				newBin, err := s.DBQueries.CreateBin(context.TODO(), database.CreateBinParams{
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

	bin, err := s.DBQueries.CreateBin(context.TODO(), database.CreateBinParams{
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

func DeleteBin(s *state.State, flags map[string]struct{}, args []string) error {
	r, v := utils.ValidateFlags(flags, "r"), utils.ValidateFlags(flags, "v")

	pathSlice := utils.ParseInputPath(args[0])

	bin, err := utils.GetLastElement(s, pathSlice)
	if err != nil {
		return err
	}

	var queue []models.Bin

	if err := utils.QueueBins(s, bin.ID, &queue); err != nil {
		return err
	}

	thisBin := models.Bin{
		Name:     bin.Name,
		ID:       bin.ID,
		ParentID: bin.ParentID,
	}

	queue = append([]models.Bin{thisBin}, queue...)
	slices.Reverse(queue)

	if r {
		for _, e := range queue {
			if v {
				fmt.Printf("deleting '%s'\n", e.Name)
			}
			err := s.DBQueries.DeleteBin(context.Background(), database.DeleteBinParams{
				Name:     e.Name,
				ParentID: e.ParentID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	}

	if len(queue) > 1 {
		return fmt.Errorf("failed to remove '%s': Bin is not empty", thisBin.Name)
	}

	parts, err := s.DBQueries.GetPartsByParent(context.Background(), queue[0].ID.UUID)
	if err != nil {
		return err
	}

	if len(parts) > 0 {
		return fmt.Errorf("failed to remove '%s': Bin is not empty", thisBin.Name)
	}

	if v {

		fmt.Printf("deleting '%s'\n", thisBin.Name)
	}

	err = s.DBQueries.DeleteBin(context.Background(), database.DeleteBinParams{
		Name:     thisBin.Name,
		ParentID: thisBin.ParentID,
	})
	if err != nil {
		return err
	}

	return nil
}

func Mv(s *state.State, flags map[string]struct{}, args []string) error {
	srcSlice := utils.ParseInputPath(args[0])
	destSlice := utils.ParseInputPath(args[1])

	fmt.Println("getting source element")
	srcElement, err := utils.GetLastElement(s, srcSlice)
	if err != nil {
		return fmt.Errorf("source path not found: %w", err)
	}

	fmt.Println("getting destination element")
	destElement, err := utils.GetLastElement(s, destSlice)
	if err != nil {
		return fmt.Errorf("source path not found: %w", err)
	}

	if srcElement.Type == "bin" && destElement.Type == "part" {
		return fmt.Errorf("cnat move a bin to a part")
	}

	destExists := destElement.ID.Valid
	elementParentID := uuid.NullUUID{Valid: false}

	if destExists {
		elementParentID = destElement.ParentID

		if destElement.ID != srcElement.ID {
			elementParentID = destElement.ID
		}
	} else {
		elementParentID = destElement.ParentID

		if srcElement.ParentID == destElement.ParentID {
			elementParentID = srcElement.ID
		}
	}

	elementName := srcElement.Name
	if !destExists && destElement.Name != "" {
		elementName = destElement.Name
	}

	if elementName != srcElement.Name {
		if srcElement.Type != destElement.Type {
			return nil
		}

		fmt.Println("updating name")

		if srcElement.Type == "bin" {
			_, err := s.DBQueries.UpdateBinName(context.Background(), database.UpdateBinNameParams{
				Name:     srcElement.Name,
				ParentID: srcElement.ParentID,
				Name_2:   elementName,
			})
			if err != nil {
				return err
			}
		}
	}

	if srcElement.ParentID != elementParentID {
		if srcElement.ID != elementParentID {
			if srcElement.Type == "part" && destElement.Type == "part" {
				return nil
			}

			fmt.Println("updating parent")

			if srcElement.Type == "part" {
				err := s.DBQueries.UpdatePartParent(context.Background(), database.UpdatePartParentParams{
					ID:       srcElement.ID.UUID,
					ParentID: elementParentID.UUID,
				})
				if err != nil {
					return err
				}
			}

			if srcElement.Type == "bin" {
				err := s.DBQueries.UpdateBinParent(context.Background(), database.UpdateBinParentParams{
					Name:       elementName,
					ParentID:   srcElement.ParentID,
					ParentID_2: elementParentID,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func UpdateBin(s *state.State, flags map[string]struct{}, args []string) error {
	v := utils.ValidateFlags(flags, "v")

	sourceSlice := utils.ParseInputPath(args[0])
	destinationSlice := utils.ParseInputPath(args[1])

	sourceParentID := uuid.NullUUID{Valid: false}
	lastBinInSource := database.Bin{}
	for _, e := range sourceSlice {
		bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
			Name:     e,
			ParentID: sourceParentID,
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
			Name:     e,
			ParentID: destinationParentID,
		})
		if err != nil {
			if i != len(destinationSlice)-1 {
				return err
			}

			bin, err := s.DBQueries.UpdateBinName(context.Background(), database.UpdateBinNameParams{
				Name:     lastBinInSource.Name,
				ParentID: destinationParentID,
				Name_2:   e,
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
		Name:       lastBinInSource.Name,
		ParentID:   lastBinInSource.ParentID,
		ParentID_2: destinationParentID,
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
	sourceSlice := utils.ParseInputPath(args[0])

	sourceParentID := uuid.NullUUID{Valid: false}
	lastBinInSource := database.Bin{}
	for _, e := range sourceSlice {
		bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
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
