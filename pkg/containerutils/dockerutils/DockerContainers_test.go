package dockerutils_test

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// Since we delete the image in this tests:
// Use a pinned ubuntu version not used in other packages might tested at the same time.
const ubuntImageName = "ubuntu:26.04"

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getRunningDockerContainerToTest(t *testing.T, implementationName string, containerName string) (containerinterfaces.Container, dockerinterfaces.Docker) {
	container, docker := getDockerContainerToTest(t, implementationName, containerName)
	err := container.Run(getCtx(), &dockeroptions.DockerRunContainerOptions{
		ImageName: "ubuntu",
		Command:   []string{"sleep", "1m"},
	})
	require.NoError(t, err)

	return container, docker
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

func TestUbuntuImageNameNotUsedInOtherPackages(t *testing.T) {
	// This is a repo test to test the repo configuration, not the implementation itself.
	// To avoid race conditions during testing we use a dedicated tagged ubuntu image.
	ctx := getCtx()

	repoRoot, err := gitutils.GetRepositoryRootPathByPath(ctx, ".")
	require.NoError(t, err)

	toCheck := ubuntImageName
	require.NotEmpty(t, toCheck)

	goFiles, err := nativefiles.ListFiles(ctx, repoRoot, &parameteroptions.ListFileOptions{
		MatchBasenamePattern: []string{`.*\.go`},
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(goFiles), 500)

	packagePath, err := pathsutils.GetAbsolutePath(".")
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(packagePath), 10)

	for _, f := range goFiles {
		if filepath.Dir(f) == packagePath {
			continue
		}

		contains, err := nativefiles.Contains(contextutils.WithSilent(ctx), f, toCheck)
		require.NoError(t, err)

		require.Falsef(t, contains, "The file '%s' contains the same image as used in the package '%s'. To avoid race conditions the '%s' package should be the only one using the image and tag '%s'.", f, packagePath, packagePath, toCheck)
	}
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
				// Use an image which is not used by other tests to avoid
				// race conditions in parallel testing when removing the image.
				const imageName = "ubuntu:24.04"
				ctx := getCtx()

				container, docker := getDockerContainerToTest(t, tt.implementationName, containername)

				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)

				if tt.enforcePullImage {
					// Delete the image so the run command is forced to perform a pull before the container can be started:
					err := docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
					require.NoError(t, err)
				} else {
					// Ensure the image is already present so no pull is needed to run the container:
					_, err := docker.PullImage(ctx, imageName)
					require.NoError(t, err)
				}

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

				err = container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
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

func TestContainers_Container_RunCommand(t *testing.T) {
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
				const imageName = ubuntImageName
				ctx := getCtx()

				container, docker := getDockerContainerToTest(t, tt.implementationName, containername)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)

				if tt.enforcePullImage {
					// Delete the image so the run command is forced to perform a pull before the container can be started:
					err := docker.RemoveImage(ctx, imageName, &dockeroptions.RemoveOptions{})
					require.NoError(t, err)
				} else {
					// Ensure the image is already present so no pull is needed to run the container:
					_, err := docker.PullImage(ctx, imageName)
					require.NoError(t, err)
				}

				// Start the container:
				err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
					Command:   []string{"sleep", "1m"},
					ImageName: imageName,
				})
				require.NoError(t, err)

				// Run another command (like docker exec) in the same container:
				output, err := container.RunCommand(ctx, &parameteroptions.RunCommandOptions{
					Command: []string{"bash", "-c", "echo hello > /world.txt; echo world"},
				})
				require.NoError(t, err)

				returnCode, err := output.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, returnCode, 0)

				stdout, err := output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, "world\n", stdout)

				// Run again a command in the same container.
				// As it is the same container we can open the file writen by the command before:
				output, err = container.RunCommand(ctx, &parameteroptions.RunCommandOptions{
					Command: []string{"cat", "/world.txt"},
				})
				require.NoError(t, err)
				stdout, err = output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, "hello\n", stdout)
			},
		)
	}
}

func Test_Container_GetHostDescription(t *testing.T) {
	tests := []struct {
		implementationName string
		containerName      string
		expected           string
	}{
		{"nativeDocker", "test-get-hostdescription", "Docker container 'test-get-hostdescription' running on host 'localhost'."},
		{"commandExectuorDockerContainer", "test-get-hostdescription", "Docker container 'test-get-hostdescription' running on host 'localhost'."},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				container, _ := getDockerContainerToTest(t, tt.implementationName, tt.containerName)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				hostDescription, err := container.GetHostDescription()
				require.NoError(t, err)
				require.EqualValues(t, tt.expected, hostDescription)
			},
		)
	}
}

func TestAdditionalEnvVarsDockerContainers(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExectuorDockerContainer"},
	}
	for _, tt := range tests {
		t.Run("env var not set "+tt.implementationName, func(t *testing.T) {
			ctx := getCtx()
			container, _ := getRunningDockerContainerToTest(t, tt.implementationName, "test-additional-env-vars-not-set")
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
			stdout, err := container.RunCommandAndGetStdoutAsString(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command: []string{"bash", "-c", "echo -en \"${MY_ENV}\""},
				},
			)
			require.NoError(t, err)
			require.Empty(t, stdout)
		})
	}

	testsValues := []struct {
		value string
	}{
		{"a"},
		{"hello"},
		{"hello world"},
		{"HELLO WORLD"},
	}

	for _, tt := range tests {
		for _, v := range testsValues {
			t.Run("env var set: "+tt.implementationName+" "+v.value, func(t *testing.T) {
				ctx := getCtx()
				container, _ := getRunningDockerContainerToTest(t, tt.implementationName, "test-additional-env-vars-not-set")
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				stdout, err := container.RunCommandAndGetStdoutAsString(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: []string{"bash", "-c", "echo -en \"${MY_ENV}\""},
						AdditionalEnvVars: map[string]string{
							"MY_ENV": v.value,
						},
					},
				)
				require.NoError(t, err)
				require.EqualValues(t, v.value, stdout)

				// Addionally the PATH variable is check to ensure it's not overwritten or absent after defining AdditionalEnvVars:
				stdout, err = container.RunCommandAndGetStdoutAsString(
					ctx,
					&parameteroptions.RunCommandOptions{
						Command: []string{"bash", "-c", "echo -en \"${PATH}\""},
						AdditionalEnvVars: map[string]string{
							"MY_ENV": v.value,
						},
					},
				)
				require.NoError(t, err)
				require.Contains(t, stdout, "/bin")
			})
		}
	}
}

func TestRunContainerWithEnvVars(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExectuorDockerContainer"},
	}
	for _, tt := range tests {
		t.Run("with env vars "+tt.implementationName, func(t *testing.T) {
			ctx := getCtx()

			const containerName = "test-run-with-env-vars"

			container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)
			err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
			require.NoError(t, err)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			err = container.Run(getCtx(), &dockeroptions.DockerRunContainerOptions{
				ImageName:         "ubuntu",
				Command:           []string{"sleep", "1m"},
				AdditionalEnvVars: map[string]string{"MY_ENV": "hello world"},
			})
			require.NoError(t, err)

			stdout, err := container.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
				Command: []string{"printenv", "MY_ENV"},
			})
			require.NoError(t, err)

			stdout = strings.TrimSpace(stdout)

			require.EqualValues(t, "hello world", stdout)
		})
	}

}

func Test_RunCommandAndGetStdoutAsIoReadCloser(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExectuorDockerContainer"},
	}
	for _, tt := range tests {
		t.Run("io.ReadCloser"+tt.implementationName, func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			const containerName = "test-run-command-and-get-stdout-as-io-read-closer"

			container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)
			err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
			require.NoError(t, err)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			err = container.Run(getCtx(), &dockeroptions.DockerRunContainerOptions{
				ImageName: "ubuntu",
				Command:   []string{"sleep", "1m"},
			})
			require.NoError(t, err)

			readCloser, err := container.RunCommandAndGetStdoutAsIoReadCloser(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command: []string{"echo", "hello", "world"},
				},
			)
			require.NoError(t, err)
			defer readCloser.Close()

			output, err := io.ReadAll(readCloser)
			require.NoError(t, err)

			require.EqualValues(t, "hello world\n", string(output))
		})
	}
}

func Test_RunCommandAndGetStdinAsIoWriteCloser(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExectuorDockerContainer"},
	}
	for _, tt := range tests {
		t.Run("io.WriteCloser"+tt.implementationName, func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			const containerName = "test-run-command-and-get-stdin-as-io-read-writer"

			container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)
			err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
			require.NoError(t, err)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			err = container.Run(getCtx(), &dockeroptions.DockerRunContainerOptions{
				ImageName: "ubuntu",
				Command:   []string{"sleep", "1m"},
			})
			require.NoError(t, err)

			writeCloser, err := container.RunCommandAndGetStdinAsIoWriteCloser(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command: []string{"tee", "/testfile"},
				},
			)
			require.NoError(t, err)
			defer writeCloser.Close()

			_, err = fmt.Fprint(writeCloser, "hello world.\n")
			require.NoError(t, err)

			err = writeCloser.Close()
			require.NoError(t, err)

			content, err := commandexecutorfile.ReadAsString(container, "/testfile")
			require.NoError(t, err)
			require.EqualValues(t, "hello world.\n", content)
		})
	}
}

func Test_Container_GetLogs(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeDocker"},
		{"commandExectuorDockerContainer"},
	}
	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt)+"_both",
			func(t *testing.T) {
				ctx := getCtx()

				const containerName = "test-get-logs"

				container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)
				err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// Run a container that produces output and exits
				err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
					ImageName:            "alpine:latest",
					Command:              []string{"sh", "-c", "echo 'stdout message' && echo 'stderr message' >&2"},
					KeepStoppedContainer: true,
				})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// Wait for container to finish using WaitUntilFinished
				err = container.WaitUntilFinished(ctx, time.Second*30)
				require.NoError(t, err)

				// Get logs
				stdout, stderr, err := container.GetLogs(ctx)
				require.NoError(t, err)

				// Verify logs contain expected messages
				stdoutStr := string(stdout)
				stderrStr := string(stderr)

				require.EqualValues(t, stdoutStr, "stdout message\n")
				require.EqualValues(t, stderrStr, "stderr message\n")
			},
		)

		t.Run(
			testutils.MustFormatAsTestname(tt)+"_stdout",
			func(t *testing.T) {
				ctx := getCtx()

				const containerName = "test-get-logs"

				container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)
				err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// Run a container that produces output and exits
				err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
					ImageName:            "alpine:latest",
					Command:              []string{"sh", "-c", "echo 'stdout message'"},
					KeepStoppedContainer: true,
				})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// Wait for container to finish using WaitUntilFinished
				err = container.WaitUntilFinished(ctx, time.Second*30)
				require.NoError(t, err)

				// Get logs
				stdout, stderr, err := container.GetLogs(ctx)
				require.NoError(t, err)

				// Verify logs contain expected messages
				stdoutStr := string(stdout)
				require.NotNil(t, stderr)
				stderrStr := string(stderr)

				require.EqualValues(t, stdoutStr, "stdout message\n")
				require.EqualValues(t, stderrStr, "")
			},
		)

		t.Run(
			testutils.MustFormatAsTestname(tt)+"_stdout",
			func(t *testing.T) {
				ctx := getCtx()

				const containerName = "test-get-logs"

				container, _ := getDockerContainerToTest(t, tt.implementationName, containerName)
				err := container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// Run a container that produces output and exits
				err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
					ImageName:            "alpine:latest",
					Command:              []string{"sh", "-c", "echo 'stderr message' >&2"},
					KeepStoppedContainer: true,
				})
				require.NoError(t, err)
				defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

				// Wait for container to finish using WaitUntilFinished
				err = container.WaitUntilFinished(ctx, time.Second*30)
				require.NoError(t, err)

				// Get logs
				stdout, stderr, err := container.GetLogs(ctx)
				require.NoError(t, err)

				// Verify logs contain expected messages
				require.NotNil(t, stdout)
				stdoutStr := string(stdout)
				stderrStr := string(stderr)

				require.EqualValues(t, stdoutStr, "")
				require.EqualValues(t, stderrStr, "stderr message\n")
			},
		)
	}
}
