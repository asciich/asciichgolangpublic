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

func TestAppendBytes(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.AppendBytes(ctx, nil, "/tmp/test.txt", []byte("hello"))
		require.Error(t, err)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.AppendBytes(ctx, commandexecutorexecoo.Exec(), "", []byte("hello"))
		require.Error(t, err)
	})

	t.Run("nil toWrite returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.AppendBytes(ctx, commandexecutorexecoo.Exec(), "/tmp/test.txt", nil)
		require.Error(t, err)
	})

	t.Run("append to existing file", func(t *testing.T) {
		ctx := getCtx()

		// Create initial file content
		tempFile := "/tmp/test_append_existing.txt"
		err := commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "initial", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		// Append more content
		err = commandexecutorfile.AppendBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte(" appended"))
		require.NoError(t, err)

		// Verify content
		content, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.Equal(t, "initial appended", content)
	})

	t.Run("append creates file if not exists", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_append_new.txt"
		// Ensure file does not exist
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.AppendBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte("new content"))
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		content, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.Equal(t, "new content", content)
	})

	t.Run("append multiple times", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_append_multiple.txt"
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.AppendBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte("first"))
		require.NoError(t, err)
		err = commandexecutorfile.AppendBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte(" second"))
		require.NoError(t, err)
		err = commandexecutorfile.AppendBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte(" third"))
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		content, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.Equal(t, "first second third", content)
	})
}

func TestAppendBytesInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-append-alpine"

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

	t.Run("append in alpine container", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_append_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.AppendBytes(ctx, container, tempFile, []byte("alpine content"))
		require.NoError(t, err)

		content, err := commandexecutorfile.ReadAsString(ctx, container, tempFile)
		require.NoError(t, err)
		require.Equal(t, "alpine content", content)

		// Append more
		err = commandexecutorfile.AppendBytes(ctx, container, tempFile, []byte(" more"))
		require.NoError(t, err)

		content, err = commandexecutorfile.ReadAsString(ctx, container, tempFile)
		require.NoError(t, err)
		require.Equal(t, "alpine content more", content)
	})
}

func TestAppendString(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.AppendString(ctx, nil, "/tmp/test.txt", "hello")
		require.Error(t, err)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.AppendString(ctx, commandexecutorexecoo.Exec(), "", "hello")
		require.Error(t, err)
	})

	t.Run("append string to file", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_append_string.txt"
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.AppendString(ctx, commandexecutorexecoo.Exec(), tempFile, "hello")
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		content, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.Equal(t, "hello", content)
	})

	t.Run("append string with newline", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_append_string_newline.txt"
		_ = commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.AppendString(ctx, commandexecutorexecoo.Exec(), tempFile, "line1\n")
		require.NoError(t, err)
		err = commandexecutorfile.AppendString(ctx, commandexecutorexecoo.Exec(), tempFile, "line2\n")
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.DeleteOptions{})

		content, err := commandexecutorfile.ReadAsString(ctx, commandexecutorexecoo.Exec(), tempFile)
		require.NoError(t, err)
		require.Equal(t, "line1\nline2\n", content)
	})
}
