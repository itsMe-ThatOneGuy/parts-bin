package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

func CreateBin(s *state.State, e models.Element) (database.Bin, error) {
	return s.DBQueries.CreateBin(context.Background(), database.CreateBinParams{
		Name:     e.Name,
		ParentID: e.ParentID,
	})
}

func CreatePart(s *state.State, e models.Element) (database.Part, error) {
	return s.DBQueries.CreatePart(context.Background(), database.CreatePartParams{
		Name:     e.Name,
		ParentID: e.ParentID.UUID,
	})
}

func GetBin(s *state.State, e models.Element) (database.Bin, error) {
	return s.DBQueries.GetBin(context.Background(), database.GetBinParams{
		Name:     e.Name,
		ParentID: e.ParentID,
	})
}

func GetPart(s *state.State, e models.Element) (database.Part, error) {
	return s.DBQueries.GetPart(context.Background(), database.GetPartParams{
		Name:     e.Name,
		Sku:      sql.NullString{Valid: true, String: e.Sku},
		ParentID: e.ParentID.UUID,
	})
}

func GetBinsByParent(s *state.State, e models.Element) ([]database.Bin, error) {
	return s.DBQueries.GetBinsByParent(context.Background(), e.ID)
}

func GetPartsByParent(s *state.State, e models.Element) ([]database.Part, error) {
	return s.DBQueries.GetPartsByParent(context.Background(), e.ID.UUID)
}

func UpdateBinName(s *state.State, name string, e models.Element) (database.Bin, error) {
	return s.DBQueries.UpdateBinName(context.Background(), database.UpdateBinNameParams{
		ID:   e.ID.UUID,
		Name: name,
	})
}

func UpdatePartName(s *state.State, name string, e models.Element) (database.Part, error) {
	return s.DBQueries.UpdatePartName(context.Background(), database.UpdatePartNameParams{
		ID:   e.ID.UUID,
		Name: name,
	})
}

func UpdateBinParent(s *state.State, parentID uuid.NullUUID, e models.Element) error {
	return s.DBQueries.UpdateBinParent(context.Background(), database.UpdateBinParentParams{
		ID:       e.ID.UUID,
		ParentID: parentID,
	})
}

func UpdatePartParent(s *state.State, parentID uuid.NullUUID, e models.Element) error {
	return s.DBQueries.UpdatePartParent(context.Background(), database.UpdatePartParentParams{
		ID:       e.ID.UUID,
		ParentID: parentID.UUID,
	})
}

func UpdateBinSku(s *state.State, sku string, e models.Element) error {
	return s.DBQueries.UpdateBinSku(context.Background(), database.UpdateBinSkuParams{
		ID: e.ID.UUID,
		Sku: sql.NullString{
			String: sku,
			Valid:  true,
		},
	})
}

func UpdatePartSku(s *state.State, sku string, e models.Element) error {
	return s.DBQueries.UpdatePartSku(context.Background(), database.UpdatePartSkuParams{
		ID: e.ID.UUID,
		Sku: sql.NullString{
			String: sku,
			Valid:  true,
		},
	})
}

func DeleteBin(s *state.State, e models.Element) error {
	return s.DBQueries.DeleteBin(context.Background(), database.DeleteBinParams{
		ID: e.ID.UUID,
	})
}

func DeletePart(s *state.State, e models.Element) error {
	return s.DBQueries.DeletePart(context.Background(), database.DeletePartParams{
		ID: e.ID.UUID,
	})
}

func DeleteManyParts(s *state.State, num int64, e models.Element) error {
	return s.DBQueries.DeleteManyParts(context.Background(), database.DeleteManyPartsParams{
		Name:     e.Name,
		ParentID: e.ParentID.UUID,
		Limit:    int32(num),
	})
}
