package commandexecutorfile_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_OpenAsWriteCloser(t *testing.T) {
	t.Run("nil commandExecutor", func(t *testing.T) {
		ctx := getCtx()

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, nil, "/tmp/test.txt", &filesoptions.WriteOptions{})
		require.Error(t, err)
		require.Nil(t, writeCloser)
	})

	t.Run("empty path", func(t *testing.T) {
		ctx := getCtx()

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.WriteOptions{})
		require.Error(t, err)
		require.Nil(t, writeCloser)
	})

	t.Run("nil options", func(t *testing.T) {
		ctx := getCtx()

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), "/tmp/test.txt", nil)
		require.Error(t, err)
		require.Nil(t, writeCloser)
	})

	t.Run("hello world", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer writeCloser.Close()

		_, err = fmt.Fprintf(writeCloser, "hello world")
		require.NoError(t, err)
		err = writeCloser.Close()
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})

	t.Run("multiple writes are concatenated", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer writeCloser.Close()

		_, err = writeCloser.Write([]byte("hello "))
		require.NoError(t, err)
		_, err = writeCloser.Write([]byte("world"))
		require.NoError(t, err)
		_, err = fmt.Fprintf(writeCloser, "%s", "!")
		require.NoError(t, err)

		err = writeCloser.Close()
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello world!", got)
	})

	t.Run("creates non existing file", func(t *testing.T) {
		ctx := getCtx()

		path := filepath.Join(t.TempDir(), "new_file.txt")

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), path, &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer writeCloser.Close()

		_, err = writeCloser.Write([]byte("created"))
		require.NoError(t, err)
		err = writeCloser.Close()
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "created", got)
	})

	t.Run("close without write results in empty file", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "previous content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		writeCloser, err := commandexecutorfile.OpenAsWriteCloser(ctx, commandexecutorexecoo.Exec(), tempFile, &filesoptions.WriteOptions{})
		require.NoError(t, err)
		err = writeCloser.Close()
		require.NoError(t, err)

		got, err := nativefiles.ReadAsBytes(ctx, tempFile)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

func Test_WriteString(t *testing.T) {
	t.Run("nil commandExecutor", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteString(ctx, nil, "/tmp/test.txt", "hello", &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("empty path", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), "", "hello", &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("write and read back hello world", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "hello world", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})

	t.Run("nil options uses defaults", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "default options", nil)
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "default options", got)
	})

	t.Run("write empty string", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		// Ensure there is previous content which must be truncated.
		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "previous content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		// []byte("") is a non-nil empty slice, so this must not fail with a nil content error.
		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "", got)
	})

	t.Run("write multiline string", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := "line1\nline2\nline3\n"
		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, content, got)
	})

	t.Run("write unicode content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := "Grüezi mitenand 🇨🇭 äöü éàè"
		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, content, got)
	})

	t.Run("overwrite existing content with shorter string", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "this is a much longer initial content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		err = commandexecutorfile.WriteString(ctx, commandexecutorexecoo.Exec(), tempFile, "short", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "short", got)
	})
}

func Test_WriteBytes(t *testing.T) {
	t.Run("nil commandExecutor", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, nil, "/tmp/test.txt", []byte("hello"), &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("empty path", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), "", []byte("hello"), &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("nil content", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), "/tmp/test.txt", nil, &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("nil options uses defaults", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte("hello world"), nil)
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})

	t.Run("write and read back hello world", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte("hello world")
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})

	t.Run("write binary content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsBytes(ctx, tempFile)
		require.NoError(t, err)
		require.EqualValues(t, content, got)
	})

	t.Run("write empty content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte{}
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsBytes(ctx, tempFile)
		require.NoError(t, err)
		require.EqualValues(t, content, got)
	})

	t.Run("write multiline content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		content := []byte("line1\nline2\nline3\n")
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "line1\nline2\nline3\n", got)
	})

	t.Run("overwrite existing content with shorter content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte("this is a much longer initial content"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, []byte("short"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, tempFile, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "short", got)
	})

	t.Run("write large content", func(t *testing.T) {
		ctx := getCtx()

		tempFile, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempFile, &filesoptions.DeleteOptions{})

		// 1 MiB of data to make sure larger pipe buffers are handled correctly.
		content := bytes.Repeat([]byte("0123456789abcdef"), 64*1024)
		err = commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), tempFile, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsBytes(ctx, tempFile)
		require.NoError(t, err)
		require.Len(t, got, len(content))
		require.True(t, bytes.Equal(content, got))
	})

	t.Run("creates non existing file", func(t *testing.T) {
		ctx := getCtx()

		path := filepath.Join(t.TempDir(), "new_file.txt")

		err := commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), path, []byte("created"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "created", got)
	})

	t.Run("creates non existing parent directories", func(t *testing.T) {
		ctx := getCtx()

		path := filepath.Join(t.TempDir(), "does", "not", "exist", "file.txt")

		err := commandexecutorfile.WriteBytes(ctx, commandexecutorexecoo.Exec(), path, []byte("hello"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := nativefiles.ReadAsString(ctx, path, &filesoptions.ReadOptions{})
		require.NoError(t, err)
		require.EqualValues(t, "hello", got)
	})
}

// Test_WriteBytesInAlpineContainer runs WriteBytes against an alpine:latest container.
// Alpine uses BusyBox instead of GNU coreutils, so this ensures only portable
// commands/options are used by WriteBytes.
func Test_WriteBytesInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Use a temporary container for testing.
	// Dedicated name to avoid collisions with other container based tests.
	const containerName = "test-commandexecutorfile-writebytes-alpine"

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

	// ensureAbsent deletes the given path inside the container and validates it's gone.
	ensureAbsent := func(t *testing.T, path string) {
		t.Helper()

		err := commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, container, path)
		require.NoError(t, err)
		require.False(t, exists)
	}

	t.Run("nil content", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_nil_content.txt"
		ensureAbsent(t, path)

		err := commandexecutorfile.WriteBytes(ctx, container, path, nil, &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("empty path", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, container, "", []byte("hello"), &filesoptions.WriteOptions{})
		require.Error(t, err)
	})

	t.Run("nil options uses defaults", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_nil_options.txt"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteBytes(ctx, container, path, []byte("hello world"), nil)
		require.NoError(t, err)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)
	})

	t.Run("creates non existing file", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_new_file.txt"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteBytes(ctx, container, path, []byte("created"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, container, path)
		require.NoError(t, err)
		require.True(t, exists)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.EqualValues(t, "created", got)
	})

	t.Run("write and read back hello world", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_hello_world.txt"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		content := []byte("hello world")
		err := commandexecutorfile.WriteBytes(ctx, container, path, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.EqualValues(t, "hello world", got)

		size, err := commandexecutorfile.GetSizeBytes(ctx, container, path)
		require.NoError(t, err)
		require.EqualValues(t, len(content), size)
	})

	t.Run("write binary content", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_binary.bin"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		// Includes NUL bytes, a newline, a quote and high bytes to catch shell quoting/encoding issues.
		content := []byte{0x00, 0x01, 0x02, '\n', '\'', '"', 0xFF, 0xFE, 0xFD, 0x00}
		err := commandexecutorfile.WriteBytes(ctx, container, path, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		size, err := commandexecutorfile.GetSizeBytes(ctx, container, path)
		require.NoError(t, err)
		require.EqualValues(t, len(content), size)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.True(t, bytes.Equal(content, []byte(got)))
	})

	t.Run("write empty content", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_empty.txt"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteBytes(ctx, container, path, []byte{}, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		exists, err := commandexecutorfile.Exists(ctx, container, path)
		require.NoError(t, err)
		require.True(t, exists)

		isEmpty, err := commandexecutorfile.IsEmptyFile(ctx, container, path)
		require.NoError(t, err)
		require.True(t, isEmpty)
	})

	t.Run("write multiline content", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_multiline.txt"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		content := []byte("line1\nline2\nline3\n")
		err := commandexecutorfile.WriteBytes(ctx, container, path, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.EqualValues(t, "line1\nline2\nline3\n", got)

		size, err := commandexecutorfile.GetSizeBytes(ctx, container, path)
		require.NoError(t, err)
		require.EqualValues(t, len(content), size)
	})

	t.Run("write unicode content", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_unicode.txt"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		content := []byte("Grüezi mitenand 🇨🇭 äöü éàè")
		err := commandexecutorfile.WriteBytes(ctx, container, path, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.EqualValues(t, string(content), got)
	})

	t.Run("overwrite existing content with shorter content", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_overwrite.txt"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteBytes(ctx, container, path, []byte("this is a much longer initial content"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		err = commandexecutorfile.WriteBytes(ctx, container, path, []byte("short"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.EqualValues(t, "short", got)

		size, err := commandexecutorfile.GetSizeBytes(ctx, container, path)
		require.NoError(t, err)
		require.EqualValues(t, len("short"), size)
	})

	t.Run("write large content", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes_large.bin"
		ensureAbsent(t, path)
		defer commandexecutorfile.Delete(ctx, container, path, &filesoptions.DeleteOptions{})

		// 1 MiB of data to make sure larger pipe buffers are handled correctly.
		content := bytes.Repeat([]byte("0123456789abcdef"), 64*1024)
		err := commandexecutorfile.WriteBytes(ctx, container, path, content, &filesoptions.WriteOptions{})
		require.NoError(t, err)

		size, err := commandexecutorfile.GetSizeBytes(ctx, container, path)
		require.NoError(t, err)
		require.EqualValues(t, len(content), size)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.Len(t, got, len(content))
		require.True(t, bytes.Equal(content, []byte(got)))
	})
	t.Run("creates non existing parent directories", func(t *testing.T) {
		ctx := getCtx()

		const path = "/tmp/writebytes/does/not/exist/file.txt"

		err := commandexecutorfile.WriteBytes(ctx, container, path, []byte("hello"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		got, err := commandexecutorfile.ReadAsString(container, path)
		require.NoError(t, err)
		require.EqualValues(t, "hello", got)
	})

	t.Run("parent path is a file returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.WriteBytes(ctx, container, "/tmp/writebytes_iam_a_file", []byte("content"), &filesoptions.WriteOptions{})
		require.NoError(t, err)

		err = commandexecutorfile.WriteBytes(ctx, container, "/tmp/writebytes_iam_a_file/file.txt", []byte("hello"), &filesoptions.WriteOptions{})
		require.Error(t, err)
	})
}
