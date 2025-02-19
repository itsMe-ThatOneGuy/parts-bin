package parts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/utils"
)

// validate the path
// create the part in the parent bin
func CreatePart(s *state.State, flags map[string]struct{}, args []string) error {
	pathSlice := utils.ParseInputPath(args[0])
	pathLen := len(pathSlice)
	part := pathSlice[len(pathSlice)-1]
	pathSlice = pathSlice[:pathLen-1]

	parentID := uuid.NullUUID{Valid: false}
	for _, e := range pathSlice {
		bin, err := s.DBQueries.GetBin(context.Background(), database.GetBinParams{
			Name:     e,
			ParentID: parentID,
		})
		if err != nil {
			return err
		}

		parentID = uuid.NullUUID{Valid: true, UUID: bin.ID}
	}

	dbPart, err := s.DBQueries.CreatePart(context.Background(), database.CreatePartParams{
		Name:     part,
		ParentID: parentID.UUID,
	})
	if err != nil {
		return err
	}

	partSku := fmt.Sprintf("%s-%d", dbPart.Name, dbPart.PartID)
	err = s.DBQueries.CreateSku(context.Background(), database.CreateSkuParams{
		PartID: dbPart.PartID,
		Sku: sql.NullString{
			String: partSku,
			Valid:  true,
		},
	})
	if err != nil {
		return nil
	}

	return nil
}
