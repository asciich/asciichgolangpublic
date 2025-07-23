package versionutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/versionutils"
)

func Test_DateVersion_GetAsString(t *testing.T) {
	t.Run("new version with current date", func(t *testing.T) {
		var version versionutils.Version
		var err error

		version = versionutils.NewCurrentDateVersion()

		versionString, err := version.GetAsString()
		require.NoError(t, err)

		require.True(t, versionutils.IsVersionString(versionString))
	})
}
