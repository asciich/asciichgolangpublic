package versionutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

func Test_Example_ReadSemanticVersionFromString(t *testing.T) {
	// Read directly as semantic version:
	version, err := versionutils.NewSmanticVersionFormString("v1.2.3")
	require.NoError(t, err)

	// Extract major, minor and patch as dedicated values
	major, minor, patch, err := version.GetMajorMinorPatch()
	require.NoError(t, err)

	// Check major, minor and patch correctly loaded.
	require.EqualValues(t, 1, major)
	require.EqualValues(t, 2, minor)
	require.EqualValues(t, 3, patch)
}
