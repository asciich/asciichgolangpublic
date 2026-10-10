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

func TestGetSha256Sum(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, nil, "/etc/hostname")
		require.Error(t, err)
		require.Empty(t, sum)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.Empty(t, sum)
	})

	t.Run("/etc/hostname returns valid sha256", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, commandexecutorexecoo.Exec(), "/etc/hostname")
		require.NoError(t, err)
		require.NotEmpty(t, sum)
		require.Len(t, sum, 64)
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, commandexecutorexecoo.Exec(), "/nonexistent/file")
		require.Error(t, err)
		require.Empty(t, sum)
	})
}

func TestGetSha256SumInAlpineContainer(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-commandexecutorfile-checksums-alpine"

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

	t.Run("/etc/hostname in alpine returns valid sha256", func(t *testing.T) {
		ctx := getCtx()

		sum, err := commandexecutorfile.GetSha256Sum(ctx, container, "/etc/hostname")
		require.NoError(t, err)
		require.NotEmpty(t, sum)
		require.Len(t, sum, 64)
	})

	t.Run("custom file in alpine returns consistent checksum", func(t *testing.T) {
		ctx := getCtx()

		tempFile := "/tmp/test_checksum_alpine.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		err := commandexecutorfile.WriteString(ctx, container, tempFile, "test content", &filesoptions.WriteOptions{})
		require.NoError(t, err)
		defer commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		sum1, err := commandexecutorfile.GetSha256Sum(ctx, container, tempFile)
		require.NoError(t, err)

		sum2, err := commandexecutorfile.GetSha256Sum(ctx, container, tempFile)
		require.NoError(t, err)

		require.Equal(t, sum1, sum2)
	})
}
