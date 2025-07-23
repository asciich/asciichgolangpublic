package kubernetesinterfaces

import (
	"context"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesparameteroptions"
)

// a generic representation of a kubernetes object like a pod, ingress, role...
type Object interface {
	CreateByYamlString(ctx context.Context, options *kubernetesparameteroptions.CreateObjectOptions) (err error)
	Delete(ctx context.Context) (err error)
	Exists(ctx context.Context) (exists bool, err error)
	GetAsYamlString() (yamlString string, err error)
	SetApiVersion(string) error
}
