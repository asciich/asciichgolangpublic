package httputilsinterfaces

import (
	"context"
)

type Server interface {
	GetPort() (port int, err error)
	GetUrl() (string, error)
	StartInBackground(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}
