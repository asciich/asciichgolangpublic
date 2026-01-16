package commandexecutordocker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockergeneric"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestGetContainerStateSatus(t *testing.T) {
	t.Run("nonexisting container", func(t *testing.T) {
		ctx := getCtx()

		commandExectuorDocker, err := commandexecutordocker.GetLocalCommandExecutorDocker()
		require.NoError(t, err)

		c, err := commandExectuorDocker.GetContainerByName("this-container-does-not-exist")
		require.NoError(t, err)

		container, ok := c.(*commandexecutordocker.CommandExecutorDockerContainer)
		require.True(t, ok)

		status, err := container.GetContainerStateStatus(ctx)
		require.Error(t, err)
		require.True(t, dockergeneric.IsErrorContainerNotFound(err))
		require.Empty(t, status)
	})
}

func Test_ContainerIsACommandExecutor(t *testing.T) {
	// It should be possible to run commands in a container in the same way as other CommandExectuors.
	// This test validates a container fullfils the CommandExecutor interface:

	var commandExecutor commandexecutorinterfaces.CommandExecutor

	commandExecutor = commandexecutordocker.NewCommandExecutorDocker()
	require.NotNil(t, commandExecutor)
}

func TestParentSet(t *testing.T) {
	docker, err := commandexecutordocker.GetLocalCommandExecutorDocker()
	require.NoError(t, err)
	require.NotNil(t, docker)

	container, err := docker.GetContainerByName("testcontainer")
	require.NoError(t, err)
	require.NotNil(t, container)

	commandExecutorContainer, ok := container.(*commandexecutordocker.CommandExecutorDockerContainer)
	require.True(t, ok)

	parent, err := commandExecutorContainer.GetParentCommandExecutorForBaseClass()
	require.NoError(t, err)
	require.NotNil(t, parent)
}
