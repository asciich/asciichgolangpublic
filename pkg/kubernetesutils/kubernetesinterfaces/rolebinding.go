package kubernetesinterfaces

import "context"

type RoleBinding interface {
	GetName() (name string, err error)
	GetNamespace() (namespace Namespace, err error)
	Exists(ctx context.Context) (exists bool, err error)
	Delete(ctx context.Context) (err error)
}
