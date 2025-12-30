package commandexecutordocker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/commandexecutordocker"
)

func Test_ImageExists(t *testing.T) {
	t.Run("image does not exist.", func(t *testing.T) {
		ctx := getCtx()

		docker, err := commandexecutordocker.GetLocalCommandExecutorDocker()
		require.NoError(t, err)

		exists, err := docker.ImageExists(ctx, "this-image-does-not-exist")
		require.NoError(t, err)
		require.False(t, exists)
	})
}
