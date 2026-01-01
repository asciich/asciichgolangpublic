package containerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

type Container interface {
	// Get the name of the container.
	GetName() (string, error)

	// Returns true if a container exists, regardless if running or not.
	Exists(ctx context.Context) (bool, error)

	IsRunning(ctx context.Context) (bool, error)

	GetHostDescription() (string, error)

	Kill(ctx context.Context) error

	// Run the container.
	Run(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) error

	// Remove a docker container.
	// Equivalent to the CLI command 'docker container remove' if docker is used.
	Remove(ctx context.Context, options *dockeroptions.RemoveOptions) error

	// Run a command as new process in the container.
	// Equivalent to the CLI command 'docker exec' if docker is used.
	RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	IsRunningOnLocalhost() (bool, error)
	RunCommandAndGetStdoutAsBytes(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]byte, error)
	RunCommandAndGetStdoutAsFloat64(ctx context.Context, options *parameteroptions.RunCommandOptions) (float64, error)
	RunCommandAndGetStdoutAsInt64(ctx context.Context, options *parameteroptions.RunCommandOptions) (int64, error)
	RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]string, error)
	RunCommandAndGetStdoutAsString(ctx context.Context, options *parameteroptions.RunCommandOptions) (string, error)
}
