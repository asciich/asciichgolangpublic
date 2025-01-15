package mapsutils

import (
	"errors"
	"sort"
)

var ErrKeyNotFound = errors.New("key not found in map")

func GetKeysOfStringMapAsSlice(input map[string]string) (keys []string) {
	if len(input) <= 0 {
		return []string{}
	}

	keys = []string{}
	for k := range input {
		keys = append(keys, k)
	}

	return keys
}

func GetKeysOfStringMapAsSliceSorted(input map[string]string) (keys []string) {
	keys = GetKeysOfStringMapAsSlice(input)
	sort.Strings(keys)
	return keys
}
