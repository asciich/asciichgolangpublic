package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestChown(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Chown(ctx, nil, "/tmp/test.txt", &parameteroptions.ChownOptions{UserName: "root"})
		require.Error(t, err)
	})

	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Chown(ctx, commandexecutorexecoo.Exec(), "", &parameteroptions.ChownOptions{UserName: "root"})
		require.Error(t, err)
	})

	t.Run("nil options returns error", func(t *testing.T) {
		ctx := getCtx()

		err := commandexecutorfile.Chown(ctx, commandexecutorexecoo.Exec(), "/tmp/test.txt", nil)
		require.Error(t, err)
	})

	t.Run("chown file to root in container", func(t *testing.T) {
		ctx := getCtx()

		docker, err := dockerutils.GetDockerOnLocalHost()
		require.NoError(t, err)

		const containerName = "test-commandexecutorfile-chown-alpine"

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

		// Install sudo in container
		_, err = container.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"apk", "add", "--no-cache", "sudo"},
		})
		require.NoError(t, err)

		tempFile := "/tmp/test_chown.txt"
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})

		// Create file first
		err = commandexecutorfile.WriteString(ctx, container, tempFile, "test content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		// Chown to root (we're already root in container, but test the function)
		err = commandexecutorfile.Chown(ctx, container, tempFile, &parameteroptions.ChownOptions{UserName: "root"})
		require.NoError(t, err)

		// Cleanup
		_ = commandexecutorfile.Delete(ctx, container, tempFile, &filesoptions.DeleteOptions{})
	})
}
