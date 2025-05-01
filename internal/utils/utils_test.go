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
	flags := make(map[string]string)
	flags["t"] = "true"

	existsTrue, _ := ValidateFlags(flags, "t")
	existsFalse, _ := ValidateFlags(flags, "f")

	if !existsTrue {
		t.Errorf("existsTrue was false when it should of been true")
	}

	if existsFalse {
		t.Errorf("existsFalse was true when it should of been false")
	}
}

func TestAbbrevName(t *testing.T) {
	name := "part"
	abbrev := AbbrevName(name)
	name1 := "test-abbrev-name-test"
	abbrev1 := AbbrevName(name1)

	if abbrev != "PAR" {
		t.Errorf("abbrev for 'test' did not match 'tes'")
	}

	if abbrev1 != "TAN" {
		t.Errorf("abbrev for 'test-abbrev-name-test' did not match 'tan'")
	}
}
