package kubernetesinterfaces

import (
	"context"
)

type DaemonSet interface {
	GetName() (name string, err error)
	GetNamespace() (namespace Namespace, err error)
	GetNamespaceName() (namespaceName string, err error)
	SetName(name string) (err error)
	Delete(ctx context.Context) (err error)
	Exists(ctx context.Context) (exists bool, err error)
}
