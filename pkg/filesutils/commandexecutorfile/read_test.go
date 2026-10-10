package commandexecutorfile_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_OpenAsReadCloser(t *testing.T) {
	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "hello world")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		readCloser, err := commandexecutorfile.OpenAsReadCloser(ctx, commandexecutorbashoo.Bash(), tempFile)
		require.NoError(t, err)
		defer readCloser.Close()

		got, err := io.ReadAll(readCloser)
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})
}

func TestReadInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-read-alpine"

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

	t.Run("read file in alpine", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_read_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, container, tempFile, "alpine content", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		content, err := commandexecutorfile.ReadAsString(ctx, container, tempFile)
		require.NoError(t, err)
		require.Equal(t, "alpine content", content)
	})
}
