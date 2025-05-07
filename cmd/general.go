package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/helptxt"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func Help(s *state.State, flags map[string]string, args []string) error {
	println(helptxt.Help)
	println(helptxt.Sig)
	println("")

	return nil
}

func Rm(s *state.State, flags map[string]string, args []string) error {
	r, _ := utils.ValidateFlags(flags, "r")
	v, _ := utils.ValidateFlags(flags, "v")
	q, qVal := utils.ValidateFlags(flags, "q")
	h, _ := utils.ValidateFlags(flags, "h")

	if h {
		println(helptxt.Rm)
		return nil
	}

	path := args[0]
	pathSlice := utils.ParseInputPath(path)

	lastElem, err := utils.GetLastElement(s, pathSlice)
	if err != nil {
		return err
	}

	if lastElem.Type == "unknown" {
		return fmt.Errorf("cannot remove '%s': No such part or bin", lastElem.Name)
	}

	if lastElem.Type == "part" {
		if q {
			num, err := strconv.ParseInt(qVal, 10, 32)
			if err != nil {
				return err
			}

			err = s.DBQueries.DeleteManyParts(context.Background(), database.DeleteManyPartsParams{
				Name:     lastElem.Name,
				ParentID: lastElem.ParentID.UUID,
				Limit:    int32(num),
			})

			return nil
		}

		err := s.DBQueries.DeletePart(context.Background(), database.DeletePartParams{
			ID: lastElem.ID.UUID,
		})
		if err != nil {
			return nil
		}

		if v {
			fmt.Printf("removed part '%s'\n", path)
		}

		return nil
	}

	thisBin := models.Bin{
		Name:     lastElem.Name,
		ID:       lastElem.ID,
		ParentID: lastElem.ParentID,
		Path:     lastElem.Path,
	}

	var queue []models.Bin

	if err := utils.QueueBins(s, path, lastElem.ID, &queue); err != nil {
		return err
	}

	queue = append([]models.Bin{thisBin}, queue...)
	slices.Reverse(queue)

	if r {
		for _, e := range queue {
			if v {
				parts, _ := s.DBQueries.GetPartsByParent(context.Background(), e.ID.UUID)
				if len(parts) >= 1 {
					for _, part := range parts {
						partName := e.Path + "/" + part.Name
						fmt.Printf("removed part: '%s'\n", partName)
					}
				}
				fmt.Printf("removed bin: '%s'\n", e.Path)
			}

			err := s.DBQueries.DeleteBin(context.Background(), database.DeleteBinParams{
				ID: e.ID.UUID,
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

	err = s.DBQueries.DeleteBin(context.Background(), database.DeleteBinParams{
		ID: lastElem.ID.UUID,
	})
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("removed bin: '%s'\n", path)
	}

	return nil
}

func Mv(s *state.State, flags map[string]string, args []string) error {
	v, _ := utils.ValidateFlags(flags, "v")
	h, _ := utils.ValidateFlags(flags, "h")

	if h {
		println(helptxt.Mv)
		return nil
	}

	srcPath := args[0]
	destPath := args[1]
	srcSlice := utils.ParseInputPath(srcPath)
	destSlice := utils.ParseInputPath(destPath)

	srcElement, err := utils.GetLastElement(s, srcSlice)
	if err != nil {
		return fmt.Errorf("source path not found: %w", err)
	}

	destElement, err := utils.GetLastElement(s, destSlice)
	if err != nil {
		return fmt.Errorf("source path not found: %w", err)
	}

	if srcElement.Type == "bin" && destElement.Type == "part" {
		return fmt.Errorf("can't move a bin to a part: %w", err)
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
		if srcElement.Type == "part" {
			part, err := s.DBQueries.UpdatePartName(context.Background(), database.UpdatePartNameParams{
				ID:   srcElement.ID.UUID,
				Name: elementName,
			})
			if err != nil {
				return err
			}

			abbrevName := utils.AbbrevName(part.Name)
			partSku := fmt.Sprintf("%s-%04d", abbrevName, part.SerialNumber.Int32)
			err = s.DBQueries.UpdatePartSku(context.Background(), database.UpdatePartSkuParams{
				ID: part.ID,
				Sku: sql.NullString{
					String: partSku,
					Valid:  true,
				},
			})

		}

		if srcElement.Type == "bin" {
			bin, err := s.DBQueries.UpdateBinName(context.Background(), database.UpdateBinNameParams{
				Name:     srcElement.Name,
				ParentID: srcElement.ParentID,
				Name_2:   elementName,
			})
			if err != nil {
				return err
			}

			abbrevName := utils.AbbrevName(bin.Name)
			binSku := fmt.Sprintf("%s-%04d", abbrevName, bin.SerialNumber.Int32)
			err = s.DBQueries.UpdateBinSku(context.Background(), database.UpdateBinSkuParams{
				ID:  bin.ID,
				Sku: sql.NullString{Valid: true, String: binSku},
			})

		}
	}

	if srcElement.ParentID != elementParentID {
		if srcElement.ID != elementParentID {
			if srcElement.Type == "part" && destElement.Type == "part" {
				return nil
			}

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

	if v {
		fmt.Printf("renamed '%s' -> '%s'\n", srcPath, destPath)
	}

	return nil
}

func Ls(s *state.State, flags map[string]string, args []string) error {
	l, _ := utils.ValidateFlags(flags, "l")
	h, _ := utils.ValidateFlags(flags, "h")

	if h {
		println(helptxt.Ls)
		return nil
	}

	srcSlice := args
	if len(args) > 0 {
		srcSlice = utils.ParseInputPath(args[0])
	}

	lastElem, err := utils.GetLastElement(s, srcSlice)
	if err != nil {
		return err
	}

	if lastElem.Type == "unknown" {
		return errors.New("not a valid bin")
	}

	if lastElem.Type == "part" {
		fmt.Printf("%s\n", lastElem.Name)
		fmt.Printf("- Sku:\t\t %s\n", lastElem.Sku)
		fmt.Printf("- ID:\t\t %v\n", lastElem.ID.UUID)
		fmt.Printf("- Created:\t %s\n", lastElem.CreatedAt)
		fmt.Printf("- Last Update:\t %s\n", lastElem.UpdatedAt)
		fmt.Print("\n")
		return nil
	}

	bins, err := s.DBQueries.GetBinsByParent(context.Background(), lastElem.ID)
	if err != nil {
		return err
	}

	parts, err := s.DBQueries.GetPartsByParent(context.Background(), lastElem.ID.UUID)
	if err != nil {
		return err
	}

	var binString string
	for i, e := range bins {
		if l {
			binString += fmt.Sprintf("%s\n", e.Name)
			binString += fmt.Sprintf("- ID:\t\t %s\n", e.ID)
			binString += fmt.Sprintf("- Created:\t %s\n", e.CreatedAt)
			binString += fmt.Sprintf("- Last Update:\t %s\n", e.UpdatedAt)
			binString += fmt.Sprint("\n")

		} else {
			binString += fmt.Sprintf("%s\t", e.Name)
			if i%5 == 4 {
				binString += "\n"
			}
		}

	}

	var partString string
	for i, e := range parts {
		if l {
			partString += fmt.Sprintf("%s\n", e.Name)
			partString += fmt.Sprintf("- Sku:\t\t %s\n", e.Sku.String)
			partString += fmt.Sprintf("- ID:\t\t %s\n", e.ID)
			partString += fmt.Sprintf("- Created:\t %s\n", e.CreatedAt)
			partString += fmt.Sprintf("- Last Update:\t %s\n", e.UpdatedAt)
			partString += fmt.Sprint("\n")
		} else {
			partString += fmt.Sprintf("%s\t", e.Name)
			if i%5 == 4 {
				partString += "\n"
			}
		}

	}

	if len(srcSlice) > 0 {
		fmt.Printf("Bin(s): %d Part(s): %d\n", len(bins), len(parts))
	} else {
		fmt.Printf("Bin(s): %d\n", len(bins))
	}

	if len(bins) > 0 {
		fmt.Println("----------")
		fmt.Println("Bins:")
		fmt.Println("----------")
		if l {
			fmt.Printf("Total: %d\n", len(bins))
			fmt.Println()
		}
		fmt.Println(binString)
	}
	if len(parts) > 0 {
		fmt.Println("----------")
		fmt.Println("Parts:")
		fmt.Println("----------")
		if l {
			fmt.Printf("Total: %d\n", len(parts))
			fmt.Println()
		}
		fmt.Println(partString)
	}

	return nil
}
