package httputils

import (
	"context"
	"errors"
)

var ErrWebServerAlreadyRunning = errors.New("web server already running")

type Server interface {
	GetPort() (port int, err error)
	StartInBackground(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}
