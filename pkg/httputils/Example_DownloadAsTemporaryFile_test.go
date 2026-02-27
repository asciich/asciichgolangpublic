package httputils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
)

// This example shows how to perform a download to a file
func Test_Example_DownloadAsTemporaryFile(t *testing.T) {
	// Preparation start...

	// Enable verbose output
	ctx := contextutils.ContextVerbose()

	// Initialize the test web server:
	const port int = 9123
	testServer, err := testwebserver.GetTestWebServer(port)
	require.NoError(t, err)
	defer testServer.Stop(ctx)
	err = testServer.StartInBackground(ctx)
	require.NoError(t, err)
	// ... preparation end.

	// To perform a GET request use:
	downloadedFile, err := httputils.DownloadAsTemporaryFile(
		// Enable progress output.
		// Since we download a very small example file we output the progress after two bytes.
		// Use 'httputils.WithDownloadProgressEveryNMBytes(ctx, 10)' to print out every 10MB.
		httpgeneric.WithDownloadProgressEveryNBytes(ctx, 2),
		// Download and output file options:
		&httpoptions.DownloadAsTemporaryFileOptions{
			RequestOptions: &httpoptions.RequestOptions{
				// Download the content of hello_world as file
				Url: "http://localhost:9123/hello_world.txt",

				// There is no need to specify the default method "GET":
				// Method: "GET",
			},
		},
	)
	// Check there was no error
	require.NoError(t, err)
	defer downloadedFile.Delete(ctx, &filesoptions.DeleteOptions{})

	// Check downloaded content
	content, err := downloadedFile.ReadAsString()
	require.NoError(t, err)
	require.EqualExportedValues(t, "hello world\n", content)
}
