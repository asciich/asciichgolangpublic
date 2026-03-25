package dockerutils_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// This example shows how to run a command in an existing container and
// get its stdout as io.ReadCloser so it can be streamed.
func Test_Example_RunCommandAndGetStdoutAsIoReadCloser(t *testing.T) {
	// use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Define the name of our test container:
	const containerName = "example-run-container-and-exec-ioreadcloser"

	// Ensure the container is absent before we start:
	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)
	exists, err := docker.ContainerExists(ctx, containerName)
	require.NoError(t, err)
	require.False(t, exists)

	// Start the container. It will automatically be exited after 30seconds (when sleep 30s ends):
	container, err := docker.RunContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name: containerName,

			// The command executed in the container:
			Command: []string{"sleep", "30s"},

			// The container image to use:
			ImageName: "ubuntu:latest",

			// Do not automatically delete a stopped container.
			KeepStoppedContainer: true,
		},
	)
	require.NoError(t, err)

	// In any case we delete the container after this test:
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	// Run a command and get the stdout as io.ReadCloser
	readCloser, err := container.RunCommandAndGetStdoutAsIoReadCloser(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"echo", "hello", "world"}, // Can also be used to "cat /a/big/file" since output is streamed in a pipe.
		},
	)
	require.NoError(t, err)
	defer readCloser.Close()

	// Read the full stdout
	output, err := io.ReadAll(readCloser)
	require.NoError(t, err)

	require.EqualValues(t, "hello world\n", string(output))
}
