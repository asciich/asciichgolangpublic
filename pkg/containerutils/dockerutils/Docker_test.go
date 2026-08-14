package dockerutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getDockerImplementationByName(implementationName string) (docker dockerinterfaces.Docker) {
	if implementationName == "commandExecutorDocker" {
		return mustutils.Must(commandexecutordocker.GetLocalCommandExecutorDocker())
	}

	if implementationName == "nativeDocker" {
		return nativedocker.NewDocker()
	}

	logging.LogFatalWithTracef("Unknown implementation name '%s'", implementationName)
	return nil
}

func TestDocker_GetHostDescription(t *testing.T) {

	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				docker := getDockerImplementationByName(tt.implementationName)

				hostDesciption, err := docker.GetHostDescription()
				require.NoError(t, err)

				require.EqualValues(t, "localhost", hostDesciption)
			},
		)
	}
}

func Test_ListContainerNames(t *testing.T) {
	list, err := dockerutils.ListContainerNames(getCtx())
	require.NoError(t, err)
	require.NotNil(t, list)
}

func Test_PullContainerImage(t *testing.T) {
	t.Run("", func(t *testing.T) {
		image, err := dockerutils.PullContainerImage(getCtx(), "alpine:latest")
		require.NoError(t, err)
		require.NotNil(t, image)

		name, err := image.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "alpine:latest", name)
	})
}

func Test_RunContainerWithEntryPoint(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const containerName = "test-entrypoint"
				err := docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)

				// Run a container with custom entrypoint that overrides the default
				// Using alpine with entrypoint set to "echo" to print arguments
				container, err := docker.RunContainer(
					ctx,
					&dockeroptions.DockerRunContainerOptions{
						Name:       containerName,
						ImageName:  "alpine:latest",
						EntryPoint: []string{"echo"},
						Command:    []string{"custom entrypoint test"},
					},
				)
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// The container should have executed "echo custom entrypoint test" and exited
				// Wait a bit for the container to finish
				time.Sleep(time.Second * 2)

				// Check that the container is not running (it should have exited after echo)
				isRunning, err := container.IsRunning(ctx)
				require.NoError(t, err)
				require.False(t, isRunning, "Container should have exited after running echo command")
			},
		)
	}
}

func Test_RunContainerWithOverwriteEntrypoint(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const imageName = "alpine:broken-entrypoint"
				image, err := docker.BuildImage(ctx, &dockeroptions.BuildContainerOptions{
					ImageNameAndTag: imageName,
					DockerfileContent: `FROM alpine:latest
Entrypoint ["this_does_not_exist.sh"]
`,
				})
				require.NoError(t, err)
				require.NotNil(t, image)

				// Test that running without overwriting entrypoint fails because
				// the broken entrypoint executable does not exist:
				t.Run("without_entrypoint_overwrite_should_fail", func(t *testing.T) {
					const containerName = "test-overwrite-entrypoint-fail"
					err := docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
					require.NoError(t, err)

					_, err = docker.RunContainer(
						ctx,
						&dockeroptions.DockerRunContainerOptions{
							Name:                 containerName,
							ImageName:            imageName,
							Command:              []string{"/bin/sh", "-c", "echo should not work && exit 0"},
							KeepStoppedContainer: true,
						},
					)
					defer docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})

					// This should fail because "this_does_not_exist.sh" is not found:
					require.Error(t, err, "Running a container with a broken entrypoint without overwriting it should fail")
				})

				// Test that running with entrypoint overwritten to empty works:
				t.Run("with_entrypoint_overwrite_should_succeed", func(t *testing.T) {
					const containerName = "test-overwrite-entrypoint-success"
					err := docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
					require.NoError(t, err)

					container, err := docker.RunContainer(
						ctx,
						&dockeroptions.DockerRunContainerOptions{
							Name:                 containerName,
							ImageName:            imageName,
							EntryPoint:           []string{}, // Explicitly overwrite entrypoint to empty
							Command:              []string{"/bin/sh", "-c", "echo overwrite test && exit 0"},
							KeepStoppedContainer: true,
						},
					)
					require.NoError(t, err)
					defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

					err = container.WaitUntilFinished(ctx, time.Second*5)
					require.NoError(t, err)

					stdout, stderr, err := container.GetLogs(ctx)
					require.NoError(t, err)

					require.EqualValues(t, "overwrite test\n", string(stdout))
					require.EqualValues(t, "", string(stderr))
				})
			},
		)
	}
}

func Test_RunCommandInTemporaryContainer(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const containerName = "test-run-command-in-temporary-container"

				// Ensure the container is absent before we start (cleanup from previous failed runs):
				err := docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
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
			},
		)
	}
}
