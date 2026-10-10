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

func TestChmod(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Chmod(ctx, nil, "/tmp/test", &filesoptions.ChmodOptions{})
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Chmod(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.ChmodOptions{})
		require.Error(t, err)
	})
	t.Run("nil options returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Chmod(ctx, commandexecutorexecoo.Exec(), "/tmp/test", nil)
		require.Error(t, err)
	})
}

func TestGetAccessPermissions(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(nil, "/tmp", false)
		require.Error(t, err)
		require.Zero(t, perms)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(commandexecutorexecoo.Exec(), "", false)
		require.Error(t, err)
		require.Zero(t, perms)
	})
	t.Run("/tmp returns valid permissions", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(commandexecutorexecoo.Exec(), "/tmp", false)
		require.NoError(t, err)
		require.NotZero(t, perms)
	})
	t.Run("nonexistent file returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissions(commandexecutorexecoo.Exec(), "/nonexistent/file", false)
		require.Error(t, err)
		require.Zero(t, perms)
	})
}

func TestGetAccessPermissionsString(t *testing.T) {
	t.Run("/etc/hostname returns valid permissions string", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissionsString(commandexecutorexecoo.Exec(), "/etc/hostname", false)
		require.NoError(t, err)
		require.NotEmpty(t, perms)
	})
	t.Run("nonexistent file returns error", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissionsString(commandexecutorexecoo.Exec(), "/nonexistent/file", false)
		require.Error(t, err)
		require.Empty(t, perms)
	})
}

func TestChmodAndGetAccessPermissionsInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-chmod-alpine"

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

	t.Run("chmod and get permissions in alpine", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_chmod_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, container, tempFile, "test", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		// Get initial permissions
		initialPerms, err := commandexecutorfile.GetAccessPermissionsString(container, tempFile, false)
		require.NoError(t, err)
		require.NotEmpty(t, initialPerms)

		err = commandexecutorfile.Chmod(ctx, container, tempFile, &filesoptions.ChmodOptions{PermissionsString: "u=rw,g=r,o=r"})
		require.NoError(t, err)

		// Verify permissions changed
		newPerms, err := commandexecutorfile.GetAccessPermissionsString(container, tempFile, false)
		require.NoError(t, err)
		require.Contains(t, newPerms, "u=rw,g=r,o=r")
	})

	t.Run("get permissions of /etc/hostname in alpine", func(t *testing.T) {
		perms, err := commandexecutorfile.GetAccessPermissionsString(container, "/etc/hostname", false)
		require.NoError(t, err)
		require.NotEmpty(t, perms)
	})
}
