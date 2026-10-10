package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func TestExists(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, nil, "/tmp")
		require.Error(t, err)
		require.False(t, exists)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, exists)
	})

	t.Run("/tmp exists", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/tmp")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("/etc exists", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/etc")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("/etc/hostname exists", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("/etc/passwd exists", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/etc/passwd")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("/dev/null exists", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/dev/null")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("nonexistent file returns false", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/tmp/this_file_does_not_exist_abc123xyz")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("nonexistent deeply nested path returns false", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/nonexistent/deeply/nested/path/file.txt")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("path with spaces does not exist", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), "/tmp/path with spaces that does not exist")
		require.NoError(t, err)
		require.False(t, exists)
	})
}

func TestExistsInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-exists-alpine"

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

	t.Run("/etc/hostname exists in alpine", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, container, "/etc/hostname")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("/bin/sh exists in alpine", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, container, "/bin/sh")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("nonexistent file returns false in alpine", func(t *testing.T) {
		ctx := getCtx()

		exists, err := commandexecutorfile.Exists(ctx, container, "/nonexistent/alpine/file")
		require.NoError(t, err)
		require.False(t, exists)
	})
}
