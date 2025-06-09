package dockerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
)

type Docker interface {
	GetContainerByName(name string) (container containerinterfaces.Container, err error)
	GetHostDescription() (hostDescription string, err error)
	ListContainers(ctx context.Context) ([]containerinterfaces.Container, error)
	ListContainerNames(ctx context.Context) ([]string, error)
}
