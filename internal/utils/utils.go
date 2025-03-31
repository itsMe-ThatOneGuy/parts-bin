package utils

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/database"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/state"
)

func ParseInputPath(path string) (pathSlice []string) {
	if len(path) < 1 {
		return []string{""}
	}
	splitSlice := strings.Split(path, "/")
	if splitSlice[0] == "" {
		splitSlice = splitSlice[1:]
	}

	return splitSlice
}

func ParseFlags(input []string, flagBool *bool) map[string]struct{} {
	flags := make(map[string]struct{})

	if strings.HasPrefix(input[1], "-") {
		*flagBool = true
		_flags := strings.Split(input[1], "")[1:]
		for _, e := range _flags {
			flags[string(e)] = struct{}{}
		}
	}

	return flags
}

func ValidateFlags(flags map[string]struct{}, key string) bool {
	_, exists := flags[key]
	return exists
}

func GetLastElement(s *state.State, path []string) (models.Element, error) {
	parentID := uuid.NullUUID{Valid: false}
	for i, e := range path {
		isLast := i == len(path)-1

		bin, err := s.DBQueries.GetBin(context.TODO(), database.GetBinParams{
			Name:     e,
			ParentID: parentID,
		})
		if err != nil {
			if isLast {
				part, err := s.DBQueries.GetPart(context.Background(), database.GetPartParams{
					Name:     e,
					ParentID: parentID.UUID,
				})
				if err == nil {
					last := models.Element{
						Type:     "part",
						Name:     part.Name,
						Sku:      part.Sku.String,
						ID:       uuid.NullUUID{Valid: true, UUID: part.ID},
						ParentID: parentID,
					}

					return last, nil
				}

				last := models.Element{
					Type:     "unknown",
					Name:     e,
					ID:       uuid.NullUUID{Valid: false},
					ParentID: parentID,
				}

				return last, nil
			}

			return models.Element{}, err
		}

		if isLast {
			last := models.Element{
				Type:     "bin",
				Name:     e,
				ID:       uuid.NullUUID{Valid: true, UUID: bin.ID},
				ParentID: parentID,
			}

			return last, nil

		}

		parentID = uuid.NullUUID{Valid: true, UUID: bin.ID}
	}

	return models.Element{}, nil
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
			ParentID: e.ParentID,
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
