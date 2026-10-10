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

func TestMove(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Move(ctx, nil, "/tmp/src.txt", "/tmp/dest.txt", &filesoptions.MoveOptions{})
		require.Error(t, err)
	})

	t.Run("empty src returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Move(ctx, commandexecutorexecoo.Exec(), "", "/tmp/dest.txt", &filesoptions.MoveOptions{})
		require.Error(t, err)
	})

	t.Run("empty dest returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Move(ctx, commandexecutorexecoo.Exec(), "/tmp/src.txt", "", &filesoptions.MoveOptions{})
		require.Error(t, err)
	})

	t.Run("move file", func(t *testing.T) {
		ctx := getCtx()

		src := "/tmp/test_move_src.txt"
		dest := "/tmp/test_move_dest.txt"

		// Clean up any existing files
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), src, &filesoptions.DeleteOptions{})
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})

		// Create source file
		err := commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), src, "test content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		// Move file
		err = commandexecutorfile.Move(ctx, commandexecutorexecoo.Exec(), src, dest, &filesoptions.MoveOptions{})
		require.NoError(t, err)

		// Verify source no longer exists
		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), src)
		require.NoError(t, err)
		require.False(t, exists)

		// Verify dest exists with correct content
		exists, err = commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), dest)
		require.NoError(t, err)
		require.True(t, exists)

		content, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), dest)
		require.NoError(t, err)
		require.Equal(t, "test content", content)

		// Cleanup
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})
	})

	t.Run("nil options uses defaults", func(t *testing.T) {
		ctx := getCtx()

		src := "/tmp/test_move_nil_options_src.txt"
		dest := "/tmp/test_move_nil_options_dest.txt"

		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), src, &filesoptions.DeleteOptions{})
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), src, "test", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.Move(ctx, commandexecutorexecoo.Exec(), src, dest, nil)
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), dest)
		require.NoError(t, err)
		require.True(t, exists)
	})
}

func TestMoveInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-move-alpine"

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

	t.Run("move file in alpine container", func(t *testing.T) {
		ctx := getCtx()

		src := "/tmp/test_move_alpine_src.txt"
		dest := "/tmp/test_move_alpine_dest.txt"

		_ = commandexecutorfile.Delete(ctx, container, src, &filesoptions.DeleteOptions{})
		_ = commandexecutorfile.Delete(ctx, container, dest, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, container, src, "alpine content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		err = commandexecutorfile.Move(ctx, container, src, dest, &filesoptions.MoveOptions{})
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, container, src)
		require.NoError(t, err)
		require.False(t, exists)

		exists, err = commandexecutorfile.Exists(ctx, container, dest)
		require.NoError(t, err)
		require.True(t, exists)

		content, err := commandexecutorfile.ReadAsString(ctx, container, dest)
		require.NoError(t, err)
		require.Equal(t, "alpine content", content)

		_ = commandexecutorfile.Delete(ctx, container, dest, &filesoptions.DeleteOptions{})
	})
}

func TestCopy(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Copy(ctx, nil, "/tmp/src.txt", "/tmp/dest.txt", &filesoptions.CopyOptions{})
		require.Error(t, err)
	})

	t.Run("empty src returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Copy(ctx, commandexecutorexecoo.Exec(), "", "/tmp/dest.txt", &filesoptions.CopyOptions{})
		require.Error(t, err)
	})

	t.Run("empty dest returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Copy(ctx, commandexecutorexecoo.Exec(), "/tmp/src.txt", "", &filesoptions.CopyOptions{})
		require.Error(t, err)
	})

	t.Run("copy file", func(t *testing.T) {
		ctx := getCtx()

		src := "/tmp/test_copy_src.txt"
		dest := "/tmp/test_copy_dest.txt"

		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), src, &filesoptions.DeleteOptions{})
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), src, "test content", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), src, &filesoptions.DeleteOptions{})
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.Copy(ctx, commandexecutorexecoo.Exec(), src, dest, &filesoptions.CopyOptions{})
		require.NoError(t, err)

		// Verify both files exist
		srcExists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), src)
		require.NoError(t, err)
		require.True(t, srcExists)

		destExists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), dest)
		require.NoError(t, err)
		require.True(t, destExists)

		// Verify content
		srcContent, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), src)
		require.NoError(t, err)

		destContent, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), dest)
		require.NoError(t, err)
		require.Equal(t, srcContent, destContent)
	})

	t.Run("nil options uses defaults", func(t *testing.T) {
		ctx := getCtx()

		src := "/tmp/test_copy_nil_options_src.txt"
		dest := "/tmp/test_copy_nil_options_dest.txt"

		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), src, &filesoptions.DeleteOptions{})
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), src, "test", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), src, &filesoptions.DeleteOptions{})
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), dest, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.Copy(ctx, commandexecutorexecoo.Exec(), src, dest, nil)
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, commandexecutorexecoo.Exec(), dest)
		require.NoError(t, err)
		require.True(t, exists)
	})
}

func TestCopyInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-copy-alpine"

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

	t.Run("copy file in alpine container", func(t *testing.T) {
		ctx := getCtx()

		src := "/tmp/test_copy_alpine_src.txt"
		dest := "/tmp/test_copy_alpine_dest.txt"

		_ = commandexecutorfile.Delete(ctx, container, src, &filesoptions.DeleteOptions{})
		_ = commandexecutorfile.Delete(ctx, container, dest, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, container, src, "alpine copy content", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, src, &filesoptions.DeleteOptions{})
		defer commandexecutorfile.Delete(ctx, container, dest, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.Copy(ctx, container, src, dest, &filesoptions.CopyOptions{})
		require.NoError(t, err)

		srcExists, err := commandexecutorfile.Exists(ctx, container, src)
		require.NoError(t, err)
		require.True(t, srcExists)

		destExists, err := commandexecutorfile.Exists(ctx, container, dest)
		require.NoError(t, err)
		require.True(t, destExists)

		srcContent, err := commandexecutorfile.ReadAsString(ctx, container, src)
		require.NoError(t, err)

		destContent, err := commandexecutorfile.ReadAsString(ctx, container, dest)
		require.NoError(t, err)
		require.Equal(t, srcContent, destContent)
	})
}
