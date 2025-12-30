package nativedocker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
)

func Test_NewContainer(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		container, err := nativedocker.NewContainer("")
		require.Error(t, err)
		require.Nil(t, container)
	})

	t.Run("simplename", func(t *testing.T) {
		container, err := nativedocker.NewContainer("simplename")
		require.NoError(t, err)
		require.NotNil(t, container)

		name, err := container.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "simplename", name)
	})
}
