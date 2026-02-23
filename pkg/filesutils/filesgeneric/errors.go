package filesgeneric

import (
	"errors"
	"os"
)

var ErrFileNotFound = errors.New("file not found")

func IsErrFileNotFound(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrFileNotFound) {
		return true
	}

	return errors.Is(err, os.ErrNotExist)
}
