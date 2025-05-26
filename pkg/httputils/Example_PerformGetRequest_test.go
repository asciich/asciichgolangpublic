package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
)

// Example how to perform a get request:
// This is the generic SendRequest returning a full response we can handle on our own.
// Hint: There are also more specify SendRequest... functions available returning directly specific values of the response as documented in the test cases below.
func Test_Example_PerformGetRequest(t *testing.T) {
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

	// To perform a GET request use:
	response, err := httputils.SendRequest(
		ctx,
		&httputils.RequestOptions{
			// Add the URL to request here:
			Url: "http://localhost:9123/hello_world.txt",

			// There is no need to specify the default method "GET":
			// Method: "GET",
		},
	)
	require.NoError(t, err)

	// To access the response body/ payload as string use:
	// Hint: There is the function SendR
	body, err := response.GetBodyAsString()
	require.NoError(t, err)

	// No we can do with the received data whatever we want:
	require.Contains(t, body, "hello world")
}

// Example how to perform a get request and directly receive the received body/paiload as string:
func Test_Example_PerformGetRequestAndGetBodyAsString(t *testing.T) {
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

	// To perform a GET request use:
	response, err := httputils.SendRequestAndGetBodyAsString(
		ctx,
		&httputils.RequestOptions{
			// Add the URL to request here:
			Url: "http://localhost:9123/hello_world.txt",

			// There is no need to specify the default method "GET":
			// Method: "GET",
		},
	)
	require.NoError(t, err)

	// No we can do with the received data whatever we want:
	require.Contains(t, response, "hello world")
}
