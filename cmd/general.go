package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/helptxt"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/repository"
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

			err = repository.DeleteManyParts(s, num, lastElem)
			if err != nil {
				return err
			}

			if v {
				fmt.Printf("removed %d parts: '%s'\n", num, path)
			}

			return nil
		}

		err := repository.DeletePart(s, lastElem)
		if err != nil {
			return nil
		}

		if v {
			fmt.Printf("removed part '%s'\n", path)
		}

		return nil
	}

	thisBin := models.Element{
		Name:     lastElem.Name,
		ID:       lastElem.ID,
		ParentID: lastElem.ParentID,
		Path:     lastElem.Path,
	}

	var queue []models.Element

	if err := utils.QueueBins(s, path, lastElem.ID, &queue); err != nil {
		return err
	}

	queue = append([]models.Element{thisBin}, queue...)
	slices.Reverse(queue)

	if r {
		for _, e := range queue {
			if v {
				parts, err := repository.GetPartsByParent(s, e)
				if err != nil {
					return err
				}

				if len(parts) >= 1 {
					for _, part := range parts {
						partName := e.Path + "/" + part.Name
						fmt.Printf("removed part: '%s'\n", partName)
					}
				}
				fmt.Printf("removed bin: '%s'\n", e.Path)
			}

			err := repository.DeleteBin(s, e)
			if err != nil {
				return err
			}

		}

		return nil
	}

	if len(queue) > 1 {
		return fmt.Errorf("failed to remove '%s': Bin is not empty", thisBin.Name)
	}

	parts, err := repository.GetPartsByParent(s, queue[0])
	if err != nil {
		return err
	}

	if len(parts) > 0 {
		return fmt.Errorf("failed to remove '%s': Bin is not empty", thisBin.Name)
	}

	err = repository.DeleteBin(s, lastElem)
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("removed bin: '%s'\n", path)
	}

	return nil
}

// needs more debugging
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
		return err
	}

	if srcElement.Type == "unknown" {
		return fmt.Errorf("source path not found")
	}

	destElement, err := utils.GetLastElement(s, destSlice)
	if err != nil {
		return err
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
			part, err := repository.UpdatePartName(s, elementName, srcElement)
			if err != nil {
				return err
			}

			abbrevName := utils.AbbrevName(elementName)
			partSku := fmt.Sprintf("%s-%04d", abbrevName, part.SerialNumber.Int32)
			err = repository.UpdatePartSku(s, partSku, srcElement)
			if err != nil {
				return err
			}

		}

		if srcElement.Type == "bin" {
			bin, err := repository.UpdateBinName(s, elementName, srcElement)
			if err != nil {
				return err
			}

			abbrevName := utils.AbbrevName(bin.Name)
			binSku := fmt.Sprintf("%s-%04d", abbrevName, bin.SerialNumber.Int32)
			err = repository.UpdateBinSku(s, binSku, srcElement)
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

			if srcElement.Type == "part" {
				err := repository.UpdatePartParent(s, elementParentID, srcElement)
				if err != nil {
					return err
				}
			}

			if srcElement.Type == "bin" {
				err := repository.UpdateBinParent(s, elementParentID, srcElement)
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
		return errors.New("not a valid path")
	}

	lastElemLongStr := ""
	if lastElem.ParentName == "" {
		lastElemLongStr += fmt.Sprintf("root/\n")
	} else {
		lastElemLongStr += fmt.Sprintf("%s/\n", lastElem.ParentName)
	}
	lastElemLongStr += fmt.Sprintf("|- %s\n", lastElem.Name)
	lastElemLongStr += fmt.Sprintf("    - Sku:\t\t %s\n", lastElem.Sku)
	lastElemLongStr += fmt.Sprintf("    - ID:\t\t %v\n", lastElem.ID.UUID)
	lastElemLongStr += fmt.Sprintf("    - Created:\t\t %s\n", lastElem.CreatedAt)
	lastElemLongStr += fmt.Sprintf("    - Last updated:\t %s\n", lastElem.CreatedAt)
	lastElemLongStr += fmt.Sprintf("    - Type:\t\t %s\n", strings.ToUpper(lastElem.Type))
	lastElemLongStr += fmt.Sprintf("    - Path:\t\t %s\n", lastElem.Path)

	if lastElem.Type == "part" {
		fmt.Println(lastElemLongStr)
		return nil
	}

	bins, err := repository.GetBinsByParent(s, lastElem)
	if err != nil {
		return err
	}

	parts, err := repository.GetPartsByParent(s, lastElem)
	if err != nil {
		return err
	}

	binString := ""
	for i, e := range bins {
		if l {
			binString += fmt.Sprintf("   |- %s\n", e.Name)
			binString += fmt.Sprintf("       - Sku:\t\t %s\n", e.Sku.String)
			binString += fmt.Sprintf("       - ID:\t\t %v\n", e.ID)
			binString += fmt.Sprintf("       - Created:\t %s\n", e.CreatedAt.Format("01-02-2006 3:4PM"))
			binString += fmt.Sprintf("       - Last updated:\t %s\n", e.UpdatedAt.Format("01-02-2006 3:4PM"))
			binString += fmt.Sprintf("       - Type:\t\t PART\n")
			binString += fmt.Sprintf("       - Path:\t\t %s/%s\n", lastElem.Path, e.Name)
		} else {
			binString += fmt.Sprintf("    %s ", e.Name)
			if i%5 == 4 {
				binString += "\n"
			}
		}

	}

	partString := ""
	for i, e := range parts {
		if l {
			partString += fmt.Sprintf("   |- %s\n", e.Name)
			partString += fmt.Sprintf("       - Sku:\t\t %s\n", e.Sku.String)
			partString += fmt.Sprintf("       - ID:\t\t %v\n", e.ID)
			partString += fmt.Sprintf("       - Created:\t %s\n", e.CreatedAt.Format("01-02-2006 3:4PM"))
			partString += fmt.Sprintf("       - Last updated:\t %s\n", e.UpdatedAt.Format("01-02-2006 3:4PM"))
			partString += fmt.Sprintf("       - Type:\t\t PART\n")
			partString += fmt.Sprintf("       - Path:\t\t %s/%s\n", lastElem.Path, e.Name)
		} else {
			partString += fmt.Sprintf("    %s", e.Name)
			if i%5 == 4 {
				partString += "\n"
			}
		}

	}

	normalStr := ""
	if l {
		if lastElem.Path == "" {
			normalStr += fmt.Sprintf("\033[4mroot/\033[0m\n")
		} else {
			normalStr += fmt.Sprintf("\033[4m/%s\033[0m\n", lastElem.Name)
			normalStr += fmt.Sprintf(" - Sku:\t\t   %s\n", lastElem.Sku)
			normalStr += fmt.Sprintf(" - ID:\t\t   %v\n", lastElem.ID.UUID)
			normalStr += fmt.Sprintf(" - Created:\t   %s\n", lastElem.CreatedAt)
			normalStr += fmt.Sprintf(" - Last updated:   %s\n", lastElem.CreatedAt)
			normalStr += fmt.Sprintf(" - Type:\t   %s\n", strings.ToUpper(lastElem.Type))
			normalStr += fmt.Sprintf(" - Path:\t   %s\n", lastElem.Path)
			if len(binString) > 0 || len(partString) > 0 {
				normalStr += "\n"
			}
		}
	} else {
		if lastElem.Path == "" {
			normalStr += fmt.Sprintf("\033[4mroot/\033[0m\n")
		} else {
			normalStr += fmt.Sprintf("\033[4m/%s\033[0m\n", lastElem.Name)
		}
	}

	if len(binString) > 0 {
		normalStr += fmt.Sprintf(" |-\033[4mBins(s)\033[0m\n")
		normalStr += binString
		normalStr += "\n"
	}
	if len(partString) > 0 {
		normalStr += fmt.Sprintf(" |-\033[4mParts(s)\033[0m\n")
		normalStr += partString
	}

	fmt.Println(normalStr)

	if lastElem.Path == "" {
		fmt.Printf("Bin(s): %d\n", len(bins))
	} else {
		fmt.Printf("Bin(s): %d Part(s): %d\n", len(bins), len(parts))
	}

	return nil
}
