package netutilserrors

import "errors"

var ErrConnectionRefused = errors.New("connection refused")

func IsConnectionRefusedError(err error) bool {
	return errors.Is(err, ErrConnectionRefused)
}
