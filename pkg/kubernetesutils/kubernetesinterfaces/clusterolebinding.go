package kubernetesinterfaces

import "context"

type ClusterRoleBinding interface {
	GetName() (name string, err error)
	Exists(ctx context.Context) (exists bool, err error)
	Delete(ctx context.Context) (err error)
}
