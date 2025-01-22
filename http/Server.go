package http

import "errors"

var ErrWebServerAlreadyRunning = errors.New("web server already running")

type Server interface {
	GetPort() (port int, err error)
	StartInBackground(verbose bool) (err error)
	Stop(verbose bool) (err error)
	MustGetPort() (port int)
	MustStartInBackground(verbose bool)
	MustStop(verbose bool)
}
