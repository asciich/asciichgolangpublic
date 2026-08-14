package dockerutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

// This example shows how to run a container with a custom ENTRYPOINT using docker.
// The ENTRYPOINT overrides the default entrypoint of the container image.
// In this example, we use alpine:latest and override its entrypoint to "echo",
// which will print the command arguments and then exit.
func Test_RunContainer_Entrypoint_Example(t *testing.T) {
	// use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Define the name of our test container:
	const containerName = "example-run-container-entrypoint"

	// Ensure the container is absent before we start:
	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)
	exists, err := docker.ContainerExists(ctx, containerName)
	require.NoError(t, err)
	require.False(t, exists)

	// Start the container with a custom entrypoint.
	// The entrypoint is set to "echo" which will print its arguments and exit.
	// The Command field provides the arguments to the entrypoint.
	container, err := docker.RunContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name: containerName,

			// The ENTRYPOINT to use (overrides the image's default entrypoint):
			EntryPoint: []string{"echo"},

			// The command arguments passed to the entrypoint:
			// This will result in: echo "Hello from custom entrypoint!"
			Command: []string{"Hello from custom entrypoint!"},

			// The container image to use:
			ImageName: "alpine:latest",

			// Do not automatically delete a stopped container.
			KeepStoppedContainer: true,
		},
	)
	require.NoError(t, err)

	// In any case we delete the container after this test:
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	err = container.WaitUntilFinished(ctx, time.Second*5)
	require.NoError(t, err)

	stdout, stderr, err := container.GetLogs(ctx)
	require.NoError(t, err)
	require.EqualValues(t, "", string(stderr))
	require.EqualValues(t, "Hello from custom entrypoint!\n", string(stdout))
}
