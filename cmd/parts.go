package cmd

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/helptxt"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/repository"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func CreatePart(s *state.State, flags map[string]string, args []string) error {
	v, _ := utils.ValidateFlags(flags, "v")
	q, qVal := utils.ValidateFlags(flags, "q")
	h, _ := utils.ValidateFlags(flags, "h")

	if h {
		println(helptxt.Mkprt)
		return nil
	}

	pathSlice := utils.ParseInputPath(args[0])
	last, err := utils.GetLastElement(s, pathSlice)
	if err != nil {
		return err
	}

	if q {
		num, err := strconv.ParseInt(qVal, 10, 32)
		if err != nil {
			return err
		}

		for i := 0; i < int(num); i++ {
			part, err := repository.CreatePart(s, last)
			if err != nil {
				return err
			}

			partElem := models.Element{
				ID: uuid.NullUUID{Valid: true, UUID: part.ID},
			}

			abbrevName := utils.AbbrevName(part.Name)
			partSku := fmt.Sprintf("%s-%04d", abbrevName, part.SerialNumber.Int32)
			err = repository.UpdatePartSku(s, partSku, partElem)
			if err != nil {
				return err
			}
		}

		if v {
			fmt.Printf("part: created part '%s' x%d\n", last.Name, int(num))
		}

		return nil
	}

	part, err := repository.CreatePart(s, last)
	if err != nil {
		return err
	}

	partElem := models.Element{
		ID: uuid.NullUUID{Valid: true, UUID: part.ID},
	}

	abbrevName := utils.AbbrevName(part.Name)
	partSku := fmt.Sprintf("%s-%04d", abbrevName, part.SerialNumber.Int32)
	err = repository.UpdatePartSku(s, partSku, partElem)
	if err != nil {
		return nil
	}

	if v {
		fmt.Printf("part: created part '%s'\n", part.Name)
	}

	return nil
}
