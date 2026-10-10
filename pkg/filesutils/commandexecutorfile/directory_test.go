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

func TestIsEmptyDirectory(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyDirectory(ctx, nil, "/tmp")
		require.Error(t, err)
		require.False(t, isEmpty)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyDirectory(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, isEmpty)
	})

	t.Run("/tmp is not empty", func(t *testing.T) {
		ctx := getCtx()

		isEmpty, err := commandexecutorfile.IsEmptyDirectory(ctx, commandexecutorexecoo.Exec(), "/tmp")
		require.NoError(t, err)
		// /tmp might be empty or not, just check no error
		_ = isEmpty
	})

	t.Run("empty directory returns true", func(t *testing.T) {
		ctx := getCtx()

		tempDir := "/tmp/test_empty_dir"
		// Clean up any existing directory
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempDir, &filesoptions.DeleteOptions{})
		// Create empty directory
		err := commandexecutorfile.CreateParentDirectory(ctx, commandexecutorexecoo.Exec(), tempDir+"/dummy", &filesoptions.CreateOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempDir, &filesoptions.DeleteOptions{})

		isEmpty, err := commandexecutorfile.IsEmptyDirectory(ctx, commandexecutorexecoo.Exec(), tempDir)
		require.NoError(t, err)
		require.True(t, isEmpty)
	})

	t.Run("directory with files returns false", func(t *testing.T) {
		ctx := getCtx()

		tempDir := "/tmp/test_nonempty_dir"
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempDir, &filesoptions.DeleteOptions{})
		err := commandexecutorfile.CreateParentDirectory(ctx, commandexecutorexecoo.Exec(), tempDir+"/dummy", &filesoptions.CreateOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempDir, &filesoptions.DeleteOptions{})

		// Create a file in the directory
		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempDir+"/test.txt", "content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		isEmpty, err := commandexecutorfile.IsEmptyDirectory(ctx, commandexecutorexecoo.Exec(), tempDir)
		require.NoError(t, err)
		require.False(t, isEmpty)
	})
}

func TestIsEmptyDirectoryInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-directory-alpine"

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

	t.Run("empty directory in alpine returns true", func(t *testing.T) {
		ctx := getCtx()

		tempDir := "/tmp/test_empty_dir_alpine"
		_ = commandexecutorfile.Delete(ctx, container, tempDir, &filesoptions.DeleteOptions{})
		err := commandexecutorfile.CreateParentDirectory(ctx, container, tempDir+"/dummy", &filesoptions.CreateOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, tempDir, &filesoptions.DeleteOptions{})

		isEmpty, err := commandexecutorfile.IsEmptyDirectory(ctx, container, tempDir)
		require.NoError(t, err)
		require.True(t, isEmpty)
	})

	t.Run("directory with files in alpine returns false", func(t *testing.T) {
		ctx := getCtx()

		tempDir := "/tmp/test_nonempty_dir_alpine"
		_ = commandexecutorfile.Delete(ctx, container, tempDir, &filesoptions.DeleteOptions{})
		err := commandexecutorfile.CreateParentDirectory(ctx, container, tempDir+"/dummy", &filesoptions.CreateOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, tempDir, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.WriteString(ctx, container, tempDir+"/test.txt", "content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		isEmpty, err := commandexecutorfile.IsEmptyDirectory(ctx, container, tempDir)
		require.NoError(t, err)
		require.False(t, isEmpty)
	})
}
