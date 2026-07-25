package kubernetesinterfaces

import (
	"context"
)

type Pod interface {
	Delete(ctx context.Context) (err error)
	Exists(ctx context.Context) (bool, error)
	GetName() (name string, err error)
	GetNamespace() (namespace Namespace, err error)
}
