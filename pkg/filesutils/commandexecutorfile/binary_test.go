package commandexecutorfile_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestIsStaticallyLinkedBinary(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, nil, "/bin/ls")
		require.Error(t, err)
		require.False(t, isStatic)
	})

	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, isStatic)
	})

	t.Run("/bin/ls is not statically linked", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "/bin/ls")
		require.NoError(t, err)
		require.False(t, isStatic)
	})

	t.Run("/bin/sh is not statically linked", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "/bin/sh")
		require.NoError(t, err)
		require.False(t, isStatic)
	})

	t.Run("nonexistent file returns false", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "/nonexistent/file")
		// The file command doesn't return an error for nonexistent files, but the output will contain "cannot open"
		// In this case we expect no error but isStatic should be false
		require.False(t, isStatic)
		// Check that the output contains error message about file not existing
		if err != nil {
			require.True(t, strings.Contains(err.Error(), "cannot open") || strings.Contains(err.Error(), "No such file"))
		}
	})
}

func TestIsStaticallyLinkedBinaryInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-binary-alpine"

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

	// Registered before 'apk add' so the container is removed even if the installation fails.
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	// alpine:latest ships neither 'file' nor a statically linked binary:
	//   - 'file' is needed by IsStaticallyLinkedBinary.
	//   - 'busybox-static' provides '/bin/busybox.static', a real statically linked binary.
	_, err = container.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"apk", "add", "--no-cache", "file", "busybox-static"},
	})
	require.NoError(t, err)

	// In alpine, /bin/ls and /bin/sh are symlinks to /bin/busybox, which is
	// dynamically linked against musl (interpreter /lib/ld-musl-*.so.1).
	t.Run("/bin/ls in alpine is not statically linked", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, container, "/bin/ls")
		require.NoError(t, err)
		require.False(t, isStatic)
	})

	t.Run("/bin/sh in alpine is not statically linked", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, container, "/bin/sh")
		require.NoError(t, err)
		require.False(t, isStatic)
	})

	t.Run("/bin/busybox.static in alpine is statically linked", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, container, "/bin/busybox.static")
		require.NoError(t, err)
		require.True(t, isStatic)
	})
}
