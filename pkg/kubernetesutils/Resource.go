package kubernetesutils

import "context"

// a generic representation of a kubernetes resource like a pod, ingress, role...
type Resource interface {
	CreateByYamlString(ctx context.Context, options *CreateResourceOptions) (err error)
	Delete(ctx context.Context) (err error)
	Exists(ctx context.Context) (exists bool, err error)
	GetAsYamlString() (yamlString string, err error)
}
