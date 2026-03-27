package dockergeneric

import "errors"

var ErrDockerContainerNotFound = errors.New("docker container not found")
var ErrDockerImageNotFound = errors.New("docker image not found")

func IsErrorContainerNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrDockerContainerNotFound)
}

func IsErrorImageNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrDockerImageNotFound)
}
