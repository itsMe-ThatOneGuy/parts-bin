package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"

	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

func CreatePart(s *state.State, flags map[string]string, args []string) error {
	v, _ := utils.ValidateFlags(flags, "v")
	q, qVal := utils.ValidateFlags(flags, "q")

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

			partSku := fmt.Sprintf("%s-%d", part.Name, part.PartID)
			err = s.DBQueries.CreateSku(context.Background(), database.CreateSkuParams{
				PartID: part.PartID,
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

	partSku := fmt.Sprintf("%s-%d", part.Name, part.PartID)
	err = s.DBQueries.CreateSku(context.Background(), database.CreateSkuParams{
		PartID: part.PartID,
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
