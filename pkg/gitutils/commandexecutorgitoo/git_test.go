package commandexecutorgitoo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
)

func TestNew(t *testing.T) {
	t.Run("nil and empty", func(t *testing.T) {
		repo, err := commandexecutorgitoo.New(nil, "")
		require.Error(t, err)
		require.Nil(t, repo)
	})

	t.Run("exec", func(t *testing.T) {
		repo, err := commandexecutorgitoo.New(commandexecutorexecoo.Exec(), "/tmp/does-not-exist")
		require.NoError(t, err)
		require.NotNil(t, repo)

		path, err := repo.GetPath()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/does-not-exist", path)

		hostDescription, err := repo.GetHostDescription()
		require.NoError(t, err)
		require.EqualValues(t, "localhost", hostDescription)

		path, hostDescription, err = repo.GetPathAndHostDescription()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/does-not-exist", path)
		require.EqualValues(t, "localhost", hostDescription)
	})
}
