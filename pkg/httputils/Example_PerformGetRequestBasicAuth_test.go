package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
)

// Example how to perform a get request with basic auth.
// The BasicAuth credentials are directly passed as option. You can also set them on the client level as shown in the example below.
//
// This is the generic SendRequest returning a full response we can handle on our own.
// Hint: There are also more specify SendRequest... functions available returning directly specific values of the response as documented in the test cases below.
func Test_Example_PerformGetRequestWithBasicAuth(t *testing.T) {
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

	// As a first step we need to get the credentials.
	// The testwebserver exposes them as a dedicated endpoints.
	// In a real world example these are the credentials you need to know...
	username, err := httputils.SendRequestAndGetBodyAsString(ctx, &httpoptions.RequestOptions{Url: "http://localhost:9123/basicauth/username"})
	require.NoError(t, err)
	password, err := httputils.SendRequestAndGetBodyAsString(ctx, &httpoptions.RequestOptions{Url: "http://localhost:9123/basicauth/password"})
	require.NoError(t, err)

	// To perform a GET request use:
	response, err := httputils.SendRequest(
		ctx,
		&httpoptions.RequestOptions{
			// Add the URL to request here:
			Url: "http://localhost:9123/basicauth/protected.txt",

			// There is no need to specify the default method "GET":
			// Method: "GET",

			BasicAuth: &httpoptions.BasicAuth{
				Username: username,
				Password: password,
			},
		},
	)
	require.NoError(t, err)

	// To access the response body/ payload as string use:
	body, err := response.GetBodyAsString()
	require.NoError(t, err)

	// No we can do with the received data whatever we want:
	require.Contains(t, body, "This is a basic auth protected message.")
}

// Example how to perform a get request with basic auth.
// The BasicAuth credentials are set on the HTTP client. You can also pass them directly in the RequestOptions as shown in the example above.
func Test_Example_PerformGetRequestWithBasicAuth_SetOnClient(t *testing.T) {
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

	// As a first step we need to get the credentials.
	// The testwebserver exposes them as a dedicated endpoints.
	// In a real world example these are the credentials you need to know...
	username, err := httputils.SendRequestAndGetBodyAsString(ctx, &httpoptions.RequestOptions{Url: "http://localhost:9123/basicauth/username"})
	require.NoError(t, err)
	password, err := httputils.SendRequestAndGetBodyAsString(ctx, &httpoptions.RequestOptions{Url: "http://localhost:9123/basicauth/password"})
	require.NoError(t, err)

	// Get the HTTP client
	client := httputils.GetNativeClient()

	// Set the basic auth credentials. They are automatically used in every request:
	err = client.SetBasicAuth(&httpoptions.BasicAuth{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)

	// To perform a GET request use:
	response, err := client.SendRequest(
		ctx,
		&httpoptions.RequestOptions{
			// Add the URL to request here:
			Url: "http://localhost:9123/basicauth/protected.txt",

			// There is no need to specify the default method "GET":
			// Method: "GET",

			// There is also no need to set the BasicAuth here as it's already set in the client.
		},
	)
	require.NoError(t, err)

	// To access the response body/ payload as string use:
	body, err := response.GetBodyAsString()
	require.NoError(t, err)

	// No we can do with the received data whatever we want:
	require.Contains(t, body, "This is a basic auth protected message.")
}
