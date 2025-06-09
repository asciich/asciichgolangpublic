package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils"
)

func Test_ListContainerNames(t *testing.T) {
	t.Run("", func(t *testing.T) {
		list, err := dockerutils.ListContainerNames(getCtx())
		require.NoError(t, err)
		require.NotNil(t, list)
	})
}
