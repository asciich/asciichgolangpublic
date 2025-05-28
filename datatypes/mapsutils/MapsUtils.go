package mapsutils

import (
	"errors"
	"sort"
)

var ErrKeyNotFound = errors.New("key not found in map")

func GetKeysOfStringMap(input map[string]string) (keys []string) {
	if len(input) <= 0 {
		return []string{}
	}

	keys = []string{}
	for k := range input {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func DeepCopyBytesMap(originalMap map[string][]byte) map[string][]byte {
	if originalMap == nil {
		return nil
	}
	
	newMap := make(map[string][]byte, len(originalMap))

	for key, value := range originalMap {
		newValue := make([]byte, len(value))
		copy(newValue, value)
		newMap[key] = newValue
	}

	return newMap
}
