package headscalegeneric

import "errors"

var ErrHeadscaleUserNotFound = errors.New("headscale user not found")

func IsErrHeadscaleUserNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrHeadscaleUserNotFound)
}
