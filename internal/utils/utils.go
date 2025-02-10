package utils

import "strings"

func ParseInputPath(s string) (last string, parent string, pathSlice []string) {
	splitSlice := strings.Split(s, "/")
	if splitSlice[0] == "" {
		splitSlice = splitSlice[1:]
	}

	lastIndex := len(splitSlice) - 1
	_last := splitSlice[lastIndex]

	if len(splitSlice) > 1 {
		parentIndex := lastIndex - 1
		_parent := splitSlice[parentIndex]
		return _last, _parent, splitSlice
	}

	return _last, "", splitSlice
}

func ValidateFlags(flags map[string]struct{}, key string) bool {
	_, exists := flags[key]
	return exists
}
