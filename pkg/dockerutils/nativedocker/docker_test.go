package nativedocker_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_GetContainerByName(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		docker := nativedocker.NewDocker()

		container, err := docker.GetContainerByName("container-name")
		require.NoError(t, err)

		name, err := container.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "container-name", name)
	})
}

func Test_GetContainerId(t *testing.T) {
	t.Run("non existing container", func(t *testing.T) {
		ctx := getCtx()

		docker := new(nativedocker.Docker)
		id, err := docker.GetContainerId(ctx, "this container does not exist")
		require.Error(t, err)
		require.True(t, dockergeneric.IsErrorContainerNotFound(err))
		require.Empty(t, id)
	})
}
