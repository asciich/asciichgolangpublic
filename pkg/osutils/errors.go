package osutils

import (
	"errors"
)

var ErrExecutableNotFound = errors.New("executable not found")

func IsExecutableNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrExecutableNotFound)
}
