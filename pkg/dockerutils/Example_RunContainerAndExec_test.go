package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// This example shows how to run a container using docker and invoke an additional command in it (like 'docker exec').
func Test_RunContainer_and_Exec_Example(t *testing.T) {
	// use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Define the name of our test container:
	const containerName = "example-run-container-and-exec"

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

	// Let's run some commands. In an untouched container /hello.txt does not exist:
	stdout, err := container.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"bash", "-c", "test -e /hello.txt && echo '/hello.txt exists.' || echo '/hello.txt does not exist.'"},
		},
	)
	require.NoError(t, err)
	require.EqualValues(t, "/hello.txt does not exist.\n", stdout)

	// Write stdin to the command:
	stdout, err = container.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			// Tee will output stdin to stdout and write it as well to /hello.txt
			Command:     []string{"tee", "/hello.txt"},
			StdinString: "hello world",
		},
	)
	require.NoError(t, err)
	require.EqualValues(t, "hello world", stdout)

	// Repeating the first command will now generate another output since the file /hello.txt exists:
	stdout, err = container.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"bash", "-c", "test -e /hello.txt && echo '/hello.txt exists.' || echo '/hello.txt does not exist.'"},
		},
	)
	require.NoError(t, err)
	require.EqualValues(t, "/hello.txt exists.\n", stdout)

	// To remvoe the container it's recommended to use defer as done in this example (right after the container was started.)
	// But anyway it can be removed at any time using:
	err = container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	exists, err = container.Exists(ctx)
	require.NoError(t,err)
	require.False(t, exists)
}
