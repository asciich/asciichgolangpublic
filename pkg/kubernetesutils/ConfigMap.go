package kubernetesutils

import "context"

type ConfigMap interface {
	Exists(ctx context.Context) (bool, error)
}
