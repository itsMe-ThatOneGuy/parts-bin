package utils

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/models"
	"github.com/itsMe-ThatOneGuy/parts-bin/internal/repository"
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

func ParseFlags(input []string, flagBool *bool) map[string]string {
	flags := make(map[string]string)

	if strings.HasPrefix(input[1], "-") {
		*flagBool = true
		raw := input[1][1:]

		for i := 0; i < len(raw); i++ {
			e := raw[i]
			if e == 'q' {
				start := i + 1
				end := start

				for end < len(raw) && raw[end] >= '0' && raw[end] <= '9' {
					end++
				}
				if end > start {
					flags["q"] = raw[start:end]
					i = end - 1
				}
			} else {
				flags[string(e)] = "true"
			}
		}
	}

	return flags
}

func ValidateFlags(flags map[string]string, key string) (bool, string) {
	value, exists := flags[key]
	return exists, value
}

func GetLastElement(s *state.State, path []string) (models.Element, error) {
	pathLen := len(path)

	last := models.Element{
		ParentID:   uuid.NullUUID{Valid: false},
		ParentName: "",
	}

	for i, e := range path {
		isLast := i == pathLen-1
		last.Name = e
		last.Sku = strings.ToUpper(e)

		bin, err := repository.GetBin(s, last)
		if err != nil {
			if isLast {
				part, err := repository.GetPart(s, last)
				if err == nil {
					last.Type = "part"
					last.Name = part.Name
					last.Sku = part.Sku.String
					last.ID = uuid.NullUUID{Valid: true, UUID: part.ID}
					last.CreatedAt = part.CreatedAt.Format("01-02-2006 3:4PM")
					last.UpdatedAt = part.UpdatedAt.Format("01-02-2006 3:4PM")
					last.Path = "/" + strings.Join(path[:], "/")

					return last, nil
				}

				last.Type = "unknown"
				last.Name = e
				last.ID = uuid.NullUUID{Valid: false}
				last.Path = "/" + strings.Join(path[:], "/")

				return last, nil
			}

			return models.Element{}, err
		}

		if isLast {
			last.Type = "bin"
			last.Name = e
			last.Sku = bin.Sku.String
			last.ID = uuid.NullUUID{Valid: true, UUID: bin.ID}
			last.CreatedAt = bin.CreatedAt.Format("01-02-2006 3:4PM")
			last.UpdatedAt = bin.UpdatedAt.Format("01-02-2006 3:4PM")
			last.Path = "/" + strings.Join(path[:], "/")

			return last, nil

		}

		last.ParentID = uuid.NullUUID{Valid: true, UUID: bin.ID}
		last.ParentName = e
	}

	return models.Element{}, nil
}

func GetChildBins(s *state.State, path string, parentID uuid.NullUUID) ([]models.Element, error) {
	parentBin := models.Element{
		ID: parentID,
	}

	bins, err := repository.GetBinsByParent(s, parentBin)
	if err != nil {
		return nil, err
	}

	binList := make([]models.Element, len(bins))
	for i, e := range bins {
		binList[i] = models.Element{
			Name:     e.Name,
			ID:       uuid.NullUUID{Valid: true, UUID: e.ID},
			ParentID: e.ParentID,
			Path:     path + "/" + e.Name,
		}
	}

	return binList, nil
}

func QueueBins(s *state.State, path string, parentID uuid.NullUUID, queue *[]models.Element) error {
	bins, err := GetChildBins(s, path, parentID)
	if err != nil {
		return err
	}

	for _, bin := range bins {
		*queue = append(*queue, bin)
		if err := QueueBins(s, path, bin.ID, queue); err != nil {
			return err
		}
	}

	return nil
}

func AbbrevName(name string) string {
	normalized := regexp.MustCompile(`\W`).ReplaceAllString(name, "_")
	words := strings.Split(normalized, "_")

	var abbrev string
	if len(words) > 1 {
		loopLen := len(words)
		if loopLen > 3 {
			loopLen = 3
		}

		for i := 0; i < loopLen; i++ {
			word := words[i]
			clnWord := regexp.MustCompile(`\d`).ReplaceAllString(word, "")
			if len(clnWord) > 0 {
				abbrev += string(clnWord[0])
			}
		}
	} else {
		cleaned := regexp.MustCompile(`\d`).ReplaceAllString(name, "")
		abbrev = cleaned[:3]
	}

	return strings.ToUpper(abbrev)
}
