package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func isDirectory(t *testing.T, commandExecutor commandexecutorinterfaces.CommandExecutor, dirPath string) bool {
	t.Helper()

	output, err := commandExecutor.RunCommand(getCtx(), &parameteroptions.RunCommandOptions{
		Command:           []string{"test", "-d", dirPath},
		AllowAllExitCodes: true,
	})
	require.NoError(t, err)

	return output.IsExitSuccess()
}

func Test_CreateParentDirectory(t *testing.T) {
	t.Run("nil commandExecutor", func(t *testing.T) {
		err := commandexecutorfile.CreateParentDirectory(getCtx(), nil, "/tmp/a/b.txt", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})

	t.Run("empty path", func(t *testing.T) {
		err := commandexecutorfile.CreateParentDirectory(getCtx(), commandexecutorexecoo.Exec(), "", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
}

func Test_CreateParentDirectoryInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-createparentdir-alpine"

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

	t.Run("creates missing nested parent directories", func(t *testing.T) {
		const filePath = "/tmp/createparent/a/b/c/file.txt"
		require.False(t, isDirectory(t, container, "/tmp/createparent/a/b/c"))

		err := commandexecutorfile.CreateParentDirectory(getCtx(), container, filePath, &filesoptions.CreateOptions{})
		require.NoError(t, err)
		require.True(t, isDirectory(t, container, "/tmp/createparent/a/b/c"))

		// The file itself must not be created:
		exists, err := commandexecutorfile.Exists(getCtx(), container, filePath)
		require.NoError(t, err)
		require.False(t, exists)

		// Writing the file must work now:
		err = commandexecutorfile.WriteString(getCtx(), container, filePath, "hello", &filesoptions.WriteOptions{})
		require.NoError(t, err)
	})

	t.Run("is idempotent", func(t *testing.T) {
		const filePath = "/tmp/createparent_idempotent/file.txt"

		for range 2 {
			err := commandexecutorfile.CreateParentDirectory(getCtx(), container, filePath, &filesoptions.CreateOptions{})
			require.NoError(t, err)
			require.True(t, isDirectory(t, container, "/tmp/createparent_idempotent"))
		}
	})

	t.Run("existing parent directory", func(t *testing.T) {
		err := commandexecutorfile.CreateParentDirectory(getCtx(), container, "/tmp/file.txt", &filesoptions.CreateOptions{})
		require.NoError(t, err)
	})

	t.Run("parent path is a file returns error", func(t *testing.T) {
		err := commandexecutorfile.WriteString(getCtx(), container, "/tmp/iam_a_file", "content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		err = commandexecutorfile.CreateParentDirectory(getCtx(), container, "/tmp/iam_a_file/file.txt", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
}
