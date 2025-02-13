package utils

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

func ParseInputPath(s string) (pathSlice []string) {
	splitSlice := strings.Split(s, "/")
	if splitSlice[0] == "" {
		splitSlice = splitSlice[1:]
	}

	return splitSlice
}

func ValidateFlags(flags map[string]struct{}, key string) bool {
	_, exists := flags[key]
	return exists
}

func GetLastBin(s *state.State, path string) (models.Bin, error) {
	pathSlice := ParseInputPath(path)

	parentID := uuid.NullUUID{Valid: false}
	for i, e := range pathSlice {
		bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
			Name:      e,
			ParentBin: parentID,
		})
		if err != nil {
			return models.Bin{}, err
		}

		if i == len(pathSlice)-1 {
			last := models.Bin{
				Name:     e,
				ID:       uuid.NullUUID{Valid: true, UUID: bin.ID},
				ParentID: parentID,
			}

			return last, nil

		}

		parentID = uuid.NullUUID{Valid: true, UUID: bin.ID}
	}

	return models.Bin{}, nil
}

func GetChildBins(s *state.State, parentID uuid.NullUUID) ([]models.Bin, error) {
	bins, err := s.DBQueries.GetBinsByParent(context.Background(), parentID)
	if err != nil {
		return nil, err
	}

	binList := make([]models.Bin, len(bins))
	for i, e := range bins {
		binList[i] = models.Bin{
			Name:     e.Name,
			ID:       uuid.NullUUID{Valid: true, UUID: e.ID},
			ParentID: e.ParentBin,
		}
	}

	return binList, nil
}

func QueueBins(s *state.State, parentID uuid.NullUUID, queue *[]models.Bin) error {
	bins, err := GetChildBins(s, parentID)
	if err != nil {
		return err
	}

	for _, bin := range bins {
		*queue = append(*queue, bin)
		if err := QueueBins(s, bin.ID, queue); err != nil {
			return err
		}
	}

	return nil
}
