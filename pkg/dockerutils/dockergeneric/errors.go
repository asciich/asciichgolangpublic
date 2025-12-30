package dockergeneric

import "errors"

var ErrDockerContainerNotFound = errors.New("docker container not found")

func IsErrorContainerNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrDockerContainerNotFound)
}