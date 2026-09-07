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
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
)

// Example how to install a file downloaded from an URL.
//
// The optional Sha256Sum ensures the downloaded file matches the expected
// checksum before it is installed.
func Test_Example_InstallFromUrl(t *testing.T) {
	// Preparation start...

	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Initialize the test web server serving the file to install:
	const port int = 9131
	testServer, err := testwebserver.GetTestWebServer(port)
	require.NoError(t, err)
	defer testServer.Stop(ctx)
	err = testServer.StartInBackground(ctx)
	require.NoError(t, err)

	// Prepare a directory to install into:
	installDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

	installPath := filepath.Join(installDir, "installed")
	// ... preparation end.

	// To install the file directly from an URL use:
	err = installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			// The URL to download and install the file from:
			SrcUrl: "http://localhost:9131/hello_world.txt",

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
	require.EqualValues(t, "hello world\n", installedContent)
}
