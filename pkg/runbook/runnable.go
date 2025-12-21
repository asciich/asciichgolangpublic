package runbook

import "context"

type Runnable interface {
	GetName() (string, error)
	GetDescription() (string, error)
	Execute(ctx context.Context) error
	Validate(ctx context.Context) error
}
