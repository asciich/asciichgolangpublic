package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestSecurelyDelete(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.SecurelyDelete(ctx, nil, "/tmp/test.txt")
		require.Error(t, err)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.SecurelyDelete(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
	})

	t.Run("securely delete existing file", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_secure_delete.txt"
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		// Create file first
		err := commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "sensitive content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		// Verify file exists
		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.True(t, exists)

		// Securely delete
		err = commandexecutorfile.SecurelyDelete(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)

		// Verify file no longer exists
		exists, err = commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("securely delete non-existing file is idempotent", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_secure_delete_nonexistent.txt"
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		// Verify file does not exist
		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.False(t, exists)

		// Securely delete should be idempotent
		err = commandexecutorfile.SecurelyDelete(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
	})
}

func TestSecurelyDeleteInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-securedelete-alpine"

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

	// Install shred in alpine (part of coreutils package)
	_, err = container.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"apk", "add", "--no-cache", "coreutils"},
	})
	require.NoError(t, err)

	t.Run("securely delete in alpine container", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_secure_delete_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, container, tempFile, "alpine sensitive", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, container, tempFile)
		require.NoError(t, err)
		require.True(t, exists)

		err = commandexecutorfile.SecurelyDelete(ctx, container, tempFile)
		require.NoError(t, err)

		exists, err = commandexecutorfile.Exists(ctx, container, tempFile)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("securely delete non-existing file in alpine is idempotent", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_secure_delete_alpine_nonexistent.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		exists, err := commandexecutorfile.Exists(ctx, container, tempFile)
		require.NoError(t, err)
		require.False(t, exists)

		err = commandexecutorfile.SecurelyDelete(ctx, container, tempFile)
		require.NoError(t, err)
	})
}
