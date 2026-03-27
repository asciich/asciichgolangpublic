package dockerutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// This example shows how to run a command in an existing container and
// get its stdin as io.WriteCloser so it can be streamed.
func Test_Example_RunCommandAndGetStdinAsIoWriteCloser(t *testing.T) {
	// use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Define the name of our test container:
	const containerName = "example-run-container-and-exec-writeCloser"

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

	// Run a command and get the stdin as io.WriteCloser
	writeCloser, err := container.RunCommandAndGetStdinAsIoWriteCloser(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"tee", "/testfile"}, // For this example we redirect the received bytes on stdin to the file '/testfile'
		},
	)
	require.NoError(t, err)
	defer writeCloser.Close()

	// write to stdin:
	_, err = fmt.Fprint(writeCloser, "hello world.\n")
	require.NoError(t, err)

	// close the writing:
	err = writeCloser.Close()
	require.NoError(t, err)

	// Read back the written file and check the content:
	content, err := commandexecutorfile.ReadAsString(container, "/testfile")
	require.NoError(t, err)
	require.EqualValues(t, "hello world.\n", content)
}
