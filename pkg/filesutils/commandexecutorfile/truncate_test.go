package commandexecutorfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
)

func TestTruncate(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Truncate(ctx, nil, "/tmp/test", 0)
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Truncate(ctx, commandexecutorexecoo.Exec(), "", 0)
		require.Error(t, err)
	})
	t.Run("negative size returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Truncate(ctx, commandexecutorexecoo.Exec(), "/tmp/test", -1)
		require.Error(t, err)
	})
}

func TestGetSizeBytes(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		size, err := commandexecutorfile.GetSizeBytes(ctx, nil, "/etc/hostname")
		require.Error(t, err)
		require.Zero(t, size)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		size, err := commandexecutorfile.GetSizeBytes(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Zero(t, size)
	})
	t.Run("/etc/hostname returns size", func(t *testing.T) {
		ctx := getCtx()
		size, err := commandexecutorfile.GetSizeBytes(ctx, commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.Greater(t, size, int64(0))
	})

	t.Run("non exsiting file", func(t *testing.T) {
		ctx := getCtx()
		size, err := commandexecutorfile.GetSizeBytes(ctx, commandexecutorexecoo.Exec(), "/does/not/exist")
		require.Error(t, err)
		require.Zero(t, size)
		require.True(t, filesgeneric.IsErrFileNotFound(err))
	})
}

// Reproduces the error seen in commandexecutoriscsi.Test_EnsureInitiorNameConfig:
//
//	stat: unrecognized option: printf
//	BusyBox v1.37.0 multi-call binary.
//
// alpine:latest uses BusyBox's stat, which doesn't support the GNU-only
// "--printf" option.
func TestGetSizeBytesAndIsEmptyFileInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Use a temporary container for testing.
	// Different name than in the iscsi test to avoid collisions when
	// packages are tested in parallel.
	const containerName = "test-commandexecutorfile-alpine"

	// Ensure the container is absent before we start
	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	// Same container setup as in Test_EnsureInitiorNameConfig
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

	// In any case we delete the container after this test
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	t.Run("GetSizeBytes of /etc/passwd returns size", func(t *testing.T) {
		ctx := getCtx()

		size, err := commandexecutorfile.GetSizeBytes(ctx, container, "/etc/passwd")
		require.NoError(t, err)
		require.Greater(t, size, int64(0))
	})

	t.Run("GetSizeBytes of /etc/hostname returns size", func(t *testing.T) {
		ctx := getCtx()

		size, err := commandexecutorfile.GetSizeBytes(ctx, container, "/etc/hostname")
		require.NoError(t, err)
		require.Greater(t, size, int64(0))
	})

	t.Run("IsEmptyFile of /etc/passwd returns false", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, container, "/etc/passwd")
		require.NoError(t, err)
		require.False(t, isEmpty)
	})

	t.Run("IsEmptyFile of nonexistent file returns error", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, container, "/tmp/this_file_does_not_exist_abc123xyz")
		require.Error(t, err)
		require.False(t, isEmpty)
	})
}

func TestIsEmptyFile(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, nil, "/tmp/some_file.txt")
		require.Error(t, err)
		require.False(t, isEmpty)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, isEmpty)
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), "/tmp/this_file_does_not_exist_abc123xyz")
		require.Error(t, err)
		require.False(t, isEmpty)
	})

	t.Run("nonexistent deeply nested path returns error", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), "/nonexistent/deeply/nested/path/file.txt")
		require.Error(t, err)
		require.False(t, isEmpty)
	})

	t.Run("empty file returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "empty.txt")
		require.NoError(t, os.WriteFile(filePath, []byte{}, 0644))

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, isEmpty)
	})

	t.Run("file with content returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "content.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("hello world\n"), 0644))

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, isEmpty)
	})

	t.Run("file with only a line break returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "newline_only.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("\n"), 0644))

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, isEmpty)
	})

	t.Run("file with only a space returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "space_only.txt")
		require.NoError(t, os.WriteFile(filePath, []byte(" "), 0644))

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, isEmpty)
	})

	t.Run("empty file in path with spaces returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "path with spaces.txt")
		require.NoError(t, os.WriteFile(filePath, []byte{}, 0644))

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, isEmpty)
	})

	t.Run("file with content in path with spaces returns false", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "path with spaces.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("content"), 0644))

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, isEmpty)
	})

	t.Run("file truncated to zero returns true", func(t *testing.T) {
		ctx := getCtx()

		filePath := filepath.Join(t.TempDir(), "to_truncate.txt")
		require.NoError(t, os.WriteFile(filePath, []byte("some content\n"), 0644))

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.False(t, isEmpty)

		err = commandexecutorfile.Truncate(ctx, commandexecutorexecoo.Exec(), filePath, 0)
		require.NoError(t, err)

		isEmpty, err = commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), filePath)
		require.NoError(t, err)
		require.True(t, isEmpty)
	})

	t.Run("/etc/passwd is not empty", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, commandexecutorexecoo.Exec(), "/etc/passwd")
		require.NoError(t, err)
		require.False(t, isEmpty)
	})
}
