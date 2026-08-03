package osutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
)

func Test_Example_WhichAll(t *testing.T) {
	// Find all occurrences of a command in PATH:
	paths, err := osutils.WhichAll("go")
	require.NoError(t, err)

	// At least one path should be found if go is installed:
	require.NotEmpty(t, paths)

	// All returned paths are absolute:
	for _, path := range paths {
		require.Contains(t, path, "/")
	}

	// Finding all occurrences of a non-existent command returns an empty slice:
	paths, err = osutils.WhichAll("nonexistentcommand12345")
	require.NoError(t, err)
	require.Empty(t, paths)

	// Passing an empty string returns an empty slice:
	paths, err = osutils.WhichAll("")
	require.Error(t, err)
	require.Empty(t, paths)
}
