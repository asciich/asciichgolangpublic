package containerinterfaces

import "context"

type Container interface {
	IsRunning(ctx context.Context) (isRunning bool, err error)
	Kill(ctx context.Context) (err error)
}
