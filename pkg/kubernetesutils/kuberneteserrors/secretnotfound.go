package kuberneteserrors

import "errors"

var ErrSecretNotFound = errors.New("kubernetes secret not found")

func IsSecretNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrSecretNotFound)
}
