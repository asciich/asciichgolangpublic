package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func TestCreateFile(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateFile(ctx, nil, "/tmp/test", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateFile(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
}

func TestCreateDirectory(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateDirectory(ctx, nil, "/tmp/test", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateDirectory(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
}

func TestCreateFileAndDirectoryInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-create-alpine"

	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	container, err := docker.RunContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:                 containerName,
			Command:              []string{"sleep", "60s"},
			ImageName:            "alpine:latest",
			KeepStoppedContainer: true,
		},
	)
	require.NoError(t, err)

	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	t.Run("create file in alpine", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_create_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.CreateFile(ctx, container, tempFile, &filesoptions.CreateOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		exists, err := commandexecutorfile.Exists(ctx, container, tempFile)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("create directory in alpine", func(t *testing.T) {
		ctx := getCtx()

		tempDir := "/tmp/test_create_dir_alpine"
		_ = commandexecutorfile.Delete(ctx, container, tempDir, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.CreateDirectory(ctx, container, tempDir, &filesoptions.CreateOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, tempDir, &filesoptions.DeleteOptions{})

		exists, err := commandexecutorfile.DirectoryExists(ctx, container, tempDir)
		require.NoError(t, err)
		require.True(t, exists)
	})
}
