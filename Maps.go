package asciichgolangpublic

import (
	"errors"
	"sort"
)

var ErrKeyNotFound = errors.New("key not found in map")

type MapsService struct{}

func Maps() (m *MapsService) {
	return NewMapsService()
}

func NewMapsService() (m *MapsService) {
	return new(MapsService)
}

func (m *MapsService) GetKeysOfStringMapAsSlice(input map[string]string) (keys []string) {
	if len(input) <= 0 {
		return []string{}
	}

	keys = []string{}
	for k := range input {
		keys = append(keys, k)
	}

	return keys
}

func (m *MapsService) GetKeysOfStringMapAsSliceSorted(input map[string]string) (keys []string) {
	keys = m.GetKeysOfStringMapAsSlice(input)
	sort.Strings(keys)
	return keys
}
