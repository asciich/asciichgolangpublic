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

func TestDelete(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Delete(ctx, nil, "/tmp/test", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("relative path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), "relative/path", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
}

func TestDeleteDirectory(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.DeleteDirectory(ctx, nil, "/tmp/test", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.DeleteDirectory(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("relative path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.DeleteDirectory(ctx, commandexecutorexecoo.Exec(), "relative/path", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
}

func TestDeleteInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-delete-alpine"

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

	t.Run("delete file in alpine", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_delete_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, container, tempFile, "test", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, container, tempFile)
		require.NoError(t, err)
		require.True(t, exists)

		err = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})
		require.NoError(t, err)

		exists, err = commandexecutorfile.Exists(ctx, container, tempFile)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("delete directory in alpine", func(t *testing.T) {
		ctx := getCtx()

		tempDir := "/tmp/test_delete_dir_alpine"
		_ = commandexecutorfile.Delete(ctx, container, tempDir, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.CreateDirectory(ctx, container, tempDir, &filesoptions.CreateOptions{})
		require.NoError(t, err)

		exists, err := commandexecutorfile.DirectoryExists(ctx, container, tempDir)
		require.NoError(t, err)
		require.True(t, exists)

		err = commandexecutorfile.DeleteDirectory(ctx, container, tempDir, &filesoptions.DeleteOptions{})
		require.NoError(t, err)

		exists, err = commandexecutorfile.DirectoryExists(ctx, container, tempDir)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("delete non-existing file is idempotent in alpine", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_delete_nonexistent_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		exists, err := commandexecutorfile.Exists(ctx, container, tempFile)
		require.NoError(t, err)
		require.False(t, exists)

		err = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})
		require.NoError(t, err)
	})
}
