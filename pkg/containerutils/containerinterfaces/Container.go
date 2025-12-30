package containerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
)

type Container interface {
	// Get the name of the container.
	GetName() (string, error)

	// Returns true if a container exists, regardless if running or not.
	Exists(ctx context.Context) (bool, error)

	IsRunning(ctx context.Context) (bool, error)
	Kill(ctx context.Context) error

	// Run the container.
	Run(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) error

	// Remove a docker container.
	// Equivalent to the CLI command 'docker container remove'
	Remove(ctx context.Context) error
}
