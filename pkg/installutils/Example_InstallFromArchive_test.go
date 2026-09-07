package installutils_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
)

// Example how to install a single file contained in a downloaded archive.
//
// The archive is downloaded from SrcUrl and the file referenced by
// SrcArchivePath is extracted and installed. Both .tar and .tar.gz archives are
// supported (the archive type is auto-detected).
//
// When a Sha256Sum is given it refers to the extracted file to install, not to
// the archive itself.
func Test_Example_InstallFromArchive(t *testing.T) {
	// Preparation start...

	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Initialize the test web server serving the archive to install from.
	// The testwebserver exposes '/hello_world.tar.gz' containing the file
	// 'hello_world.txt' with the content "hello world\n".
	const port int = 9132
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

	// The checksum refers to the extracted file to install, not the archive:
	sha256sum := checksumutils.GetSha256SumFromString("hello world\n")
	// ... preparation end.

	// To install a file contained in a downloaded (.tar or .tar.gz) archive use:
	err = installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			// The URL to download the archive from:
			SrcUrl: "http://localhost:9132/hello_world.tar.gz",

			// The path inside the archive to the file to install:
			SrcArchivePath: "hello_world.txt",

			// The path the extracted file should be installed to:
			InstallPath: installPath,

			// The access permissions the installed file should get:
			Mode: "u=rwx",

			// The expected checksum of the extracted file (optional):
			Sha256Sum: sha256sum,
		},
	)
	require.NoError(t, err)

	// Now the extracted file is installed and we can read its content:
	installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
	require.NoError(t, err)
	require.EqualValues(t, "hello world\n", installedContent)
}
