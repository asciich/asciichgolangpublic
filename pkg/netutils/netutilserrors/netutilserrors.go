package netutilserrors

import "errors"

var ErrConnectionRefused = errors.New("connection refused")
var ErrNoRouteToHost = errors.New("no route to host")

func IsConnectionRefusedError(err error) bool {
	return errors.Is(err, ErrConnectionRefused)
}

func IsNoRouteToHostError(err error) bool {
	return errors.Is(err, ErrNoRouteToHost)
}
