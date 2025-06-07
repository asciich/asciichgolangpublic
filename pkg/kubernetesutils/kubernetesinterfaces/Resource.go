package kubernetesinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
)

// a generic representation of a kubernetes resource like a pod, ingress, role...
type Resource interface {
	CreateByYamlString(ctx context.Context, options *kubernetesparameteroptions.CreateResourceOptions) (err error)
	Delete(ctx context.Context) (err error)
	Exists(ctx context.Context) (exists bool, err error)
	GetAsYamlString() (yamlString string, err error)
}
