package versionutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

func Test_DateVersion_GetAsString(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var version versionutils.Version
		var err error

		version = versionutils.NewCurrentDateVersion()

		versionString, err := version.GetAsString()
		require.NoError(t, err)

		require.True(t, versionutils.IsVersionString(versionString))
	})
}
