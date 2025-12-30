package commandexecutordocker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
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
