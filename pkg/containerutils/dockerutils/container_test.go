package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestContainers_Container_StopAndStart(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExectuorDockerContainer"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const containerName = "test-stop-and-start-container"
				const imageName = "alpine:latest"

				container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)

				err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
					ImageName: imageName,
					Command:   []string{"sh", "-c", "sleep 1m"},
					// Do not automatically delete a stopped container.
					// Otherwise the container could not be started again after stopping:
					KeepStoppedContainer: true,
				})
				require.NoError(t, err)

				// The container should be running now:
				isRunning, err := container.IsRunning(ctx)
				require.NoError(t, err)
				require.True(t, isRunning)

				// Stop the running container:
				err = container.Stop(ctx)
				require.NoError(t, err)

				isRunning, err = container.IsRunning(ctx)
				require.NoError(t, err)
				require.False(t, isRunning)

				// The stopped container still exists:
				exists, err := container.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				// Stopping a stopped container is idempotent:
				err = container.Stop(ctx)
				require.NoError(t, err)

				isRunning, err = container.IsRunning(ctx)
				require.NoError(t, err)
				require.False(t, isRunning)

				// Start the stopped container again:
				err = container.Start(ctx)
				require.NoError(t, err)

				isRunning, err = container.IsRunning(ctx)
				require.NoError(t, err)
				require.True(t, isRunning)

				// Starting a running container is idempotent:
				err = container.Start(ctx)
				require.NoError(t, err)

				isRunning, err = container.IsRunning(ctx)
				require.NoError(t, err)
				require.True(t, isRunning)
			},
		)
	}
}

func TestContainers_Container_StartNotExistingContainer(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExectuorDockerContainer"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const containerName = "test-start-not-existing-container"

				container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)

				err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// Starting a container which does not exist returns an error:
				err = container.Start(ctx)
				require.Error(t, err)
			},
		)
	}
}
