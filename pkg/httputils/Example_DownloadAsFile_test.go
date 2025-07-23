package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils/httputilsparameteroptions"
)

// This example shows how to perform a download to a file
func Test_Example_DownloadAsFile(t *testing.T) {
	// Preparation start...

	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Initialize the test web server:
	const port int = 9123
	testServer, err := httputils.GetTestWebServer(port)
	require.NoError(t, err)
	defer testServer.Stop(ctx)
	err = testServer.StartInBackground(ctx)
	require.NoError(t, err)
	// ... preparation end.

	outputPath := "/tmp/example-download"

	// To perform a GET request use:
	downloadedFile, err := httputils.DownloadAsFile(
		// Enable progress output.
		// Since we download a very small example file we output the progress after two bytes.
		// Use 'httputils.WithDownloadProgressEveryNMBytes(ctx, 10)' to print out every 10MB.
		httputils.WithDownloadProgressEveryNBytes(ctx, 2),
		// Download and output file options:
		&httputilsparameteroptions.DownloadAsFileOptions{
			RequestOptions: &httputilsparameteroptions.RequestOptions{
				// Download the content of hello_world as file
				Url: "http://localhost:9123/hello_world.txt",

				// There is no need to specify the default method "GET":
				// Method: "GET",
			},
			OutputPath:        outputPath,
			OverwriteExisting: true,
		},
	)
	// Check there was no error
	require.NoError(t, err)

	// Check the downloaded file has the same path as requested:
	downloadedPath, err := downloadedFile.GetPath()
	require.NoError(t, err)
	require.EqualValues(t, outputPath, downloadedPath)

	// Check downloaded content
	content, err := downloadedFile.ReadAsString()
	require.NoError(t, err)
	require.EqualExportedValues(t, "hello world\n", content)
}
