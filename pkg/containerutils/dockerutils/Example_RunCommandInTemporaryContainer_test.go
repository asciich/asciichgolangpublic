package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

// This example shows how to run a command in a temporary Docker container.
// The container is created automatically, executes the command, and is deleted afterwards.
// This is similar to `docker run --rm` but provides programmatic access to the output.
func Test_Example_RunCommandInTemporaryContainer(t *testing.T) {
	// use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Define the name of our temporary container:
	const containerName = "example-run-command-in-temporary-container"

	// Ensure the container is absent before we start (cleanup from previous failed runs):
	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	// Run a command in a temporary container. The container will be automatically
	// created, started, waited for completion, logs fetched, and then deleted.
	output, err := docker.RunCommandInTemporaryContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:      containerName,
			ImageName: "alpine:latest",
			Command:   []string{"echo", "hello", "from", "temporary", "container"},
		},
	)
	require.NoError(t, err)

	// Read the stdout of the executed command:
	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)
	require.EqualValues(t, "hello from temporary container\n", stdout)

	// Verify the container was cleaned up:
	exists, err := docker.ContainerExists(ctx, containerName)
	require.NoError(t, err)
	require.False(t, exists, "Temporary container should have been removed")
}
