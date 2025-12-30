package dockerutils_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getDockerContainerToTest(t *testing.T, implementationName string, containerName string) (containerinterfaces.Container, dockerinterfaces.Docker) {
	if implementationName == "commandExectuorDockerContainer" {
		docker, err := commandexecutordocker.GetLocalCommandExecutorDocker()
		require.NoError(t, err)

		container, err := docker.GetContainerByName(containerName)
		require.NoError(t, err)
		return container, docker
	}
	if implementationName == "nativeDocker" {
		docker := nativedocker.NewDocker()
		container, err := docker.GetContainerByName(containerName)
		require.NoError(t, err)
		return container, docker
	}

	logging.LogFatalWithTracef("Unkown implementaion name: '%s'", implementationName)

	return nil, nil
}

func TestContainers_Container_Run(t *testing.T) {
	tests := []struct {
		enforcePullImage   bool
		implementationName string
	}{
		{false, "nativeDocker"},
		{true, "nativeDocker"},
		{false, "commandExectuorDockerContainer"},
		{true, "commandExectuorDockerContainer"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const containername = "test-run-container"
				const imageName = "ubuntu:latest"
				ctx := getCtx()

				container, docker := getDockerContainerToTest(t, tt.implementationName, "containername-"+strings.ToLower(tt.implementationName))
				if tt.enforcePullImage {
					// Delete the image so the run command is forced to perform a pull before the container can be started:
					err := docker.RemoveImage(ctx, imageName)
					require.NoError(t, err)
				} else {
					// Ensure the image is already present so no pull is needed to run the container:
					_, err := docker.PullImage(ctx, imageName)
					require.NoError(t, err)
				}

				defer container.Remove(ctx)
				err := container.Remove(ctx)
				require.NoError(t, err)

				// Test a deleted container does not exist:
				exists, err := container.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				// Test a deleted container is not considered running:
				isRunning, err := container.IsRunning(ctx)
				require.NoError(t, err)
				require.False(t, isRunning)

				err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
					ImageName:            imageName,
					Command:              []string{"sleep", "10s"},
					KeepStoppedContainer: true,
				})
				require.NoError(t, err)
				defer container.Kill(ctx)

				exists, err = container.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				isRunning, err = container.IsRunning(ctx)
				require.NoError(t, err)
				require.True(t, isRunning)

				err = container.Kill(ctx)
				require.NoError(t, err)

				exists, err = container.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				isRunning, err = container.IsRunning(ctx)
				require.NoError(t, err)
				require.False(t, isRunning)

				err = container.Remove(ctx)
				require.NoError(t, err)

				exists, err = container.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				isRunning, err = container.IsRunning(ctx)
				require.NoError(t, err)
				require.False(t, isRunning)
			},
		)
	}
}
