package kubernetesinterfaces

import (
	"context"
)

type Deployment interface {
	Create(ctx context.Context) (err error)
	Delete(ctx context.Context) (err error)
	Exists(ctx context.Context) (bool, error)
	GetName() (name string, err error)
	GetNamespace() (namespace string, err error)
	GetReplicas() (replicas int32, err error)
	GetAvailableReplicas() (availableReplicas int32, err error)
	GetDesiredReplicas() (desiredReplicas int32, err error)
	WaitUntilAvailable(ctx context.Context, timeoutSeconds int) (err error)
	WaitUntilDeleted(ctx context.Context, timeoutSeconds int) (err error)
}
