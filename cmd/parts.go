package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/helptxt"

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
			part, err := s.DBQueries.CreatePart(context.Background(), database.CreatePartParams{
				Name:     last.Name,
				ParentID: last.ParentID.UUID,
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

		if v {
			fmt.Printf("part: created part '%s' x%d\n", last.Name, int(num))
		}

		return nil
	}

	part, err := s.DBQueries.CreatePart(context.Background(), database.CreatePartParams{
		Name:     last.Name,
		ParentID: last.ParentID.UUID,
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
	if err != nil {
		return nil
	}

	if v {
		fmt.Printf("part: created part '%s'\n", part.Name)
	}

	return nil
}
