package kubernetesinterfaces

import "context"

type CronJob interface {
	Exists(ctx context.Context) (bool, error)
}
