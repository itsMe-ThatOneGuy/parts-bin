package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/helptxt"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/repository"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func CreateBin(s *state.State, flags map[string]string, args []string) error {
	p, _ := utils.ValidateFlags(flags, "p")
	v, _ := utils.ValidateFlags(flags, "v")
	h, _ := utils.ValidateFlags(flags, "h")

	if h {
		println(helptxt.Mkbin)
		return nil
	}

	pathSlice := utils.ParseInputPath(args[0])

	last := pathSlice[len(pathSlice)-1]
	currentBin := models.Element{
		Name: last,
	}

	if len(pathSlice) > 1 {
		currentBin.ParentID = uuid.NullUUID{Valid: false}

		for i, e := range pathSlice {
			lastEle := i == len(pathSlice)-1

			currentBin.Name = e

			bin, err := repository.GetBin(s, currentBin)
			if err != nil {
				if !lastEle && !p {
					return fmt.Errorf("cannot create bin '%s': No such parent bin", e)
				}

				newBin, err := repository.CreateBin(s, currentBin)
				if err != nil {
					return fmt.Errorf("issue creating bin '%s': %v", e, err)
				}

				currentBin.ID = uuid.NullUUID{Valid: true, UUID: newBin.ID}

				abbrevName := utils.AbbrevName(newBin.Name)
				binSku := fmt.Sprintf("%s-%04d", abbrevName, newBin.SerialNumber.Int32)
				err = repository.UpdateBinSku(s, binSku, currentBin)
				if err != nil {
					return err
				}

				currentBin.ParentID = currentBin.ID

				if v {
					fmt.Printf("created bin '%s'\n", newBin.Name)
				}

			} else {
				currentBin.ParentID = uuid.NullUUID{Valid: true, UUID: bin.ID}

				if v {
					fmt.Printf("bin: cannot create bin '%s': bin exists\n", bin.Name)
				}

			}
		}

		return nil
	}

	bin, err := repository.CreateBin(s, currentBin)
	if err != nil {
		return err
	}

	currentBin.ID = uuid.NullUUID{Valid: true, UUID: bin.ID}

	abbrevName := utils.AbbrevName(last)
	binSku := fmt.Sprintf("%s-%04d", abbrevName, bin.SerialNumber.Int32)
	err = repository.UpdateBinSku(s, binSku, currentBin)
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("bin: created bin '%s'\n", bin.Name)
	}

	return nil
}
