package http

import "errors"

var ErrWebServerAlreadyRunning = errors.New("web server already running")

type Server interface {
	StartInBackground(verbose bool) (err error)
	Stop(verbose bool) (err error)
	MustStartInBackground(verbose bool)
	MustStop(verbose bool)
}
