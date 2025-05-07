package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/helptxt"
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

				abbrevName := utils.AbbrevName(newBin.Name)
				binSku := fmt.Sprintf("%s-%04d", abbrevName, newBin.SerialNumber.Int32)
				err = s.DBQueries.UpdateBinSku(context.Background(), database.UpdateBinSkuParams{
					ID:  newBin.ID,
					Sku: sql.NullString{Valid: true, String: binSku},
				})

				parentID = uuid.NullUUID{Valid: true, UUID: newBin.ID}

				if v {
					fmt.Printf("created bin '%s'\n", newBin.Name)
				}

			} else {
				parentID = uuid.NullUUID{Valid: true, UUID: bin.ID}

				if v {
					fmt.Printf("bin: cannot create bin '%s': bin exists\n", bin.Name)
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

	abbrevName := utils.AbbrevName(last)
	binSku := fmt.Sprintf("%s-%04d", abbrevName, bin.SerialNumber.Int32)
	err = s.DBQueries.UpdateBinSku(context.Background(), database.UpdateBinSkuParams{
		ID:  bin.ID,
		Sku: sql.NullString{Valid: true, String: binSku},
	})

	if v {
		fmt.Printf("bin: created bin '%s'\n", bin.Name)
	}

	return nil
}
