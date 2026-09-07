package installutils_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
)

// Example how to install a file from a local source path.
//
// The convenience function installutils.Install delegates to the native
// implementation and is the simplest way to install a file on the localhost.
func Test_Example_InstallFromFile(t *testing.T) {
	// Preparation start...

	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Create a local source file to install.
	// In a real world example this is the file you already have and want to install.
	sourceFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "example file content\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, sourceFile, &filesoptions.DeleteOptions{})

	// Prepare a directory to install into:
	installDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

	installPath := filepath.Join(installDir, "installed")
	// ... preparation end.

	// To install the file from a local source path use:
	err = installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			// The local source file to install:
			SrcPath: sourceFile,

			// The path the file should be installed to:
			InstallPath: installPath,

			// The access permissions the installed file should get:
			Mode: "u=rwx",
		},
	)
	require.NoError(t, err)

	// Now the file is installed and we can read its content:
	installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.EqualValues(t, "example file content\n", installedContent)
}
