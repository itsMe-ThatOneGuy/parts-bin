package utils

import (
	"reflect"
	"testing"
)

func TestParseInputPath(t *testing.T) {
	inputPath := "/dir/dir1/dir2"

	expected := []string{"dir", "dir1", "dir2"}
	parsed := ParseInputPath(inputPath)

	if !reflect.DeepEqual(expected, parsed) {
		t.Errorf("parsed slice does not equal expected slice")
	}
}

func TestValidateFlags(t *testing.T) {
	flags := make(map[string]struct{})
	flags["t"] = struct{}{}

	existsTrue := ValidateFlags(flags, "t")
	existsFalse := ValidateFlags(flags, "f")

	if !existsTrue {
		t.Errorf("existsTrue was false when it should of been true")
	}

	if existsFalse {
		t.Errorf("existsFalse was true when it should of been false")
	}
}
