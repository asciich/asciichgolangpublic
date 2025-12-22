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

// Simple example how to perform a POST request.
func Test_Example_PostRequest(t *testing.T) {
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

	// This payload will be send to the webserver using a POST request.
	payload := "This is the payload for the post example."

	// Use a get request to receive the YAML data:
	response, err := httputils.SendRequest(
		ctx,
		&httpoptions.RequestOptions{
			// Add the URL to request here:
			Url: "http://localhost:9123/return_post_payload",

			// The POST method must be explicitly specified:
			Method: "POST",

			// Add the data to send.
			Data: []byte(payload),
		},
	)
	require.NoError(t, err)

	// The 'return_post_payload' returns the received data byte by byte.
	// So we can check if the correct data was send to the server and received:
	bodyString, err := response.GetBodyAsString()
	require.NoError(t, err)
	require.EqualValues(t, payload, bodyString)
}
