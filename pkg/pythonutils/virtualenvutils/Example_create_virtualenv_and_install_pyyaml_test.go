package virtualenvutils_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/pythonutils/virtualenvutils"
)

func Test_CreateVirtualenvAndInstallPyyaml(t *testing.T) {
	// Use your context here:
	ctx := context.Background()

	// Enable verbose output
	ctx = contextutils.WithVerbose(ctx)

	// Get a temporary path to install the virtualenv:
	tempDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})
	vEnvPath := filepath.Join(tempDir, "ve")

	// The virtualenv does not exist yet:
	exists, err := nativefiles.Exists(ctx, vEnvPath)
	require.NoError(t, err)
	require.False(t, exists)

	// Create the virtualenv with pyyaml installed
	ve, err := virtualenvutils.CreateVirtualEnv(ctx, &virtualenvutils.CreateVirtualenvOptions{
		Path:     vEnvPath,
		Packages: []string{"pyyaml"},
	})
	require.NoError(t, err)

	// The virtualenv now exists:
	exists, err = nativefiles.Exists(ctx, vEnvPath)
	require.NoError(t, err)
	require.True(t, exists)

	// Check pyyaml installed
	pyyamlInstalled, err := ve.IsPackageInstalled(ctx, "pyyaml")
	require.NoError(t, err)
	require.True(t, pyyamlInstalled)
}
