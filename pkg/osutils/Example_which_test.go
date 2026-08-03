package osutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
)

func Test_Example_Which(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())
	_ = ctx // ctx is available for verbose logging

	// Find the path of a common command like "go":
	path, err := osutils.Which("go")
	require.NoError(t, err)
	require.NotEmpty(t, path)

	// The returned path is absolute:
	require.Contains(t, path, "/")

	// Finding a non-existent command returns ErrExecutableNotFound:
	_, err = osutils.Which("nonexistentcommand12345")
	require.ErrorIs(t, err, osutils.ErrExecutableNotFound)

	// Passing an empty string returns an error:
	_, err = osutils.Which("")
	require.Error(t, err)
}
