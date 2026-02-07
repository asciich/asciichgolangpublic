package nativedocker_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func Test_NewContainer(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		container, err := nativedocker.NewContainer("")
		require.Error(t, err)
		require.Nil(t, container)
	})

	t.Run("simplename", func(t *testing.T) {
		container, err := nativedocker.NewContainer("simplename")
		require.NoError(t, err)
		require.NotNil(t, container)

		name, err := container.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "simplename", name)
	})
}

func Test_ContainerIsACommandExecutor(t *testing.T) {
	// It should be possible to run commands in a container in the same way as other CommandExectuors.
	// This test validates a container fullfils the CommandExecutor interface:

	var commandExecutor commandexecutorinterfaces.CommandExecutor
	var err error

	commandExecutor, err = nativedocker.NewContainer("containername")
	require.NoError(t, err)
	require.NotNil(t, commandExecutor)
}

func Test_RunCommand(t *testing.T) {
	var tests = []struct {
		cmd string
	}{
		{cmd: "printenv"},
		{cmd: "/usr/bin/printenv"},
	}

	for _, tt := range tests {
		t.Run("path variable using "+tt.cmd, func(t *testing.T) {
			ctx := getCtx()
			const containerName = "test-nativedocker-path-variable"
			docker := nativedocker.NewDocker()
			require.NoError(t, docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true}))

			container, err := docker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
				ImageName: "ubuntu",
				Name:      containerName,
				Command:   []string{"sleep", "1m"},
			})
			require.NoError(t, err)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			stdout, err := container.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
				Command: []string{tt.cmd, "PATH"},
			})
			require.NoError(t, err)
			require.Contains(t, stdout, "/bin")
		})
	}
}

func Test_RunContainerWithVolume(t *testing.T) {
	t.Run("mounted file", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world\n")
		require.NoError(t, err)

		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		const containerName = "test-mount-file"
		container, err := nativedocker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
			KeepStoppedContainer: false,
			Name:                 containerName,
			ImageName:            "ubuntu",
			Command:              []string{"sleep", "1m"},
			Mounts:               []string{tempFile + ":/etc/examplefile.txt"},
		})
		require.NoError(t, err)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		stdout, err := container.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"cat", "/etc/examplefile.txt"},
		})
		require.NoError(t, err)
		require.EqualValues(t, "hello world\n", stdout)
	})

	t.Run("mounted dir", func(t *testing.T) {
		ctx := getCtx()

		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

		err = nativefiles.WriteString(ctx, filepath.Join(tempDir, "examplefile.txt"), "hello world\n")
		require.NoError(t, err)

		err = nativefiles.WriteString(ctx, filepath.Join(tempDir, "examplefile2.txt"), "another file")
		require.NoError(t, err)

		const containerName = "test-mount-directory"
		container, err := nativedocker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
			KeepStoppedContainer: false,
			Name:                 containerName,
			ImageName:            "ubuntu",
			Command:              []string{"sleep", "1m"},
			Mounts:               []string{tempDir + ":/etc/directory"},
		})
		require.NoError(t, err)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		stdout, err := container.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"cat", "/etc/directory/examplefile.txt"},
		})
		require.NoError(t, err)
		require.EqualValues(t, "hello world\n", stdout)

		stdout, err = container.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"cat", "/etc/directory/examplefile2.txt"},
		})
		require.NoError(t, err)
		require.EqualValues(t, "another file", stdout)
	})
}

func Test_DeleteContainerTwice(t *testing.T) {
	ctx := getCtx()

	const containerName = "test-delete-container-twice"
	container, err := nativedocker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
		KeepStoppedContainer: false,
		Name:                 containerName,
		ImageName:            "ubuntu",
		Command:              []string{"sleep", "1m"},
	})
	require.NoError(t, err)
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	exists, err := container.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// delete the container the first time:
	ctxDelete := contextutils.WithChangeIndicator(ctx)
	err = container.Remove(ctxDelete, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)
	require.True(t, contextutils.IsChanged(ctxDelete))

	// delete the container the second time to check idempotence:
	ctxDelete = contextutils.WithChangeIndicator(ctx)
	err = container.Remove(ctxDelete, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)
	require.False(t, contextutils.IsChanged(ctxDelete))
}
