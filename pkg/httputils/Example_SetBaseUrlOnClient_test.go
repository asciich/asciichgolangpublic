package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpnativeclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
)

// This example demonstrates how to configure an HTTP client with a predefined base URL.
//
// Setting the base URL eliminates the need to construct the full endpoint URL
// (protocol, host, and port) for every subsequent request, allowing the user
// to only specify the request path.
func Test_Example_SetBaseUrlOnClient_test(t *testing.T) {
	// Preparation start...

	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Initialize the test web server:
	const port int = 9123
	testServer, err := testwebserver.GetTestWebServer(port)
	require.NoError(t, err)
	defer testServer.Stop(ctx)
	err = testServer.StartInBackground(ctx)
	require.NoError(t, err)
	// ... preparation end.

	// Get the HTTP client
	client := httpnativeclientoo.NewNativeClient()

	// Set the base URL and port
	err = client.SetBaseUrl("http://localhost")
	require.NoError(t, err)

	err = client.SetPort(port)
	require.NoError(t, err)

	// To perform a GET request use:
	response, err := client.SendRequest(
		ctx,
		&httpoptions.RequestOptions{
			// We only need to specify the path
			Path: "hello_world.txt",
		},
	)
	require.NoError(t, err)

	// To access the response body/ payload as string use:
	body, err := response.GetBodyAsString()
	require.NoError(t, err)

	// Now we can do with the received data whatever we want:
	require.EqualValues(t, "hello world\n", body)

	// To perform another request we use a shorter version.
	// Anyway we only have to adjust the path:
	body2, err := client.SendRequestAndGetBodyAsString(ctx, &httpoptions.RequestOptions{
		Path: "hello_world2.txt",
	})
	require.NoError(t, err)
	require.EqualValues(t, "hello world2\n", body2)
}
