package dockerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
)

type Docker interface {
	ContainerExists(ctx context.Context, name string) (bool, error)

	GetContainerByName(name string) (containerinterfaces.Container, error)
	GetImageByName(name string) (containerinterfaces.Image, error)
	GetHostDescription() (string, error)

	ImageExists(ctx context.Context, name string) (bool, error)

	KillContainerByName(ctx context.Context, name string) error

	ListContainers(ctx context.Context) ([]containerinterfaces.Container, error)
	ListContainerNames(ctx context.Context) ([]string, error)

	PullImage(ctx context.Context, imageName string) (containerinterfaces.Image, error)

	RemoveImage(ctx context.Context, imageName string) error
	RunContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (containerinterfaces.Container, error)

	RemoveContainer(ctx context.Context, name string, options *dockeroptions.RemoveOptions) error
}
