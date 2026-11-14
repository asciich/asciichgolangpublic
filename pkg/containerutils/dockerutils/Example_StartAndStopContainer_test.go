package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

// This example shows how to stop and start a container using docker.
// A stopped container is not removed: It can be started again afterwards.
func Test_StartAndStopContainer_Example(t *testing.T) {
	// use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Define the name of our test container:
	const containerName = "example-start-and-stop-container"

	// Ensure the container is absent before we start:
	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)
	exists, err := docker.ContainerExists(ctx, containerName)
	require.NoError(t, err)
	require.False(t, exists)

	// Run the container in the background.
	// The container runs 'sleep 1m' and therefore stays running:
	container, err := docker.RunContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name: containerName,

			// The command to run in the container:
			Command: []string{"sleep", "1m"},

			// The container image to use:
			ImageName: "alpine:latest",

			// Do not automatically delete a stopped container.
			// Otherwise the container could not be started again after stopping:
			KeepStoppedContainer: true,
		},
	)
	require.NoError(t, err)

	// In any case we remove the container after this test:
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	// The container should be running now:
	isRunning, err := container.IsRunning(ctx)
	require.NoError(t, err)
	require.True(t, isRunning)

	// Stop the container.
	// The container is stopped but not removed:
	err = container.Stop(ctx)
	require.NoError(t, err)

	isRunning, err = container.IsRunning(ctx)
	require.NoError(t, err)
	require.False(t, isRunning)

	exists, err = container.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Start the stopped container again:
	err = container.Start(ctx)
	require.NoError(t, err)

	isRunning, err = container.IsRunning(ctx)
	require.NoError(t, err)
	require.True(t, isRunning)
}
