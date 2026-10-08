package commandexecutoruserutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/userutils/commandexecutoruserutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestIsRunningAsRoot_NilCommandExecutor(t *testing.T) {
	ctx := getCtx()

	_, err := commandexecutoruserutils.IsRunningAsRoot(ctx, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "commandExecutor")
}

func TestIsRunningAsRoot_DockerContainerAsRoot(t *testing.T) {
	ctx := getCtx()

	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	const containerName = "test-isrunningasroot-docker-root"

	// Clean up any existing container
	err = docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	// Run a container as root (default)
	container, err := docker.RunContainer(
		ctx,
		&dockeroptions.DockerRunContainerOptions{
			Name:      containerName,
			ImageName: "alpine:latest",
			Command:   []string{"sleep", "1m"},
		},
	)
	require.NoError(t, err)
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	// Test IsRunningAsRoot - should return true for root user
	isRoot, err := commandexecutoruserutils.IsRunningAsRoot(ctx, container)
	require.NoError(t, err)
	require.True(t, isRoot, "Docker container running as root should return true")
}
