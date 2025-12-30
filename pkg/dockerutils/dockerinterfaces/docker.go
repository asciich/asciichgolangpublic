package dockerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
)

type Docker interface {
	GetContainerByName(name string) (containerinterfaces.Container, error)
	GetHostDescription() (string, error)

	KillContainerByName(ctx context.Context, name string) error

	ListContainers(ctx context.Context) ([]containerinterfaces.Container, error)
	ListContainerNames(ctx context.Context) ([]string, error)

	RunContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (containerinterfaces.Container, error)
}
