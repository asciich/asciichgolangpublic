package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils/httputilsimplementationindependend"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils/httputilsparameteroptions"
)

// This example shows how a 404 response of the webserver is handled.
func Test_Example_PerformGetRequest404(t *testing.T) {
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
		&httputilsparameteroptions.RequestOptions{
			// This URL does not exits:
			Url: "http://localhost:9123/this-page-does-not-exist.txt",

			// There is no need to specify the default method "GET":
			// Method: "GET",
		},
	)
	// If the return value is not ok an error is returned:
	require.Error(t, err)
	require.ErrorIs(t, err, httputilsimplementationindependend.ErrUnexpectedStatusCode)

	// Even the the return value was not ok the response is returned:
	require.NotNil(t, response)

	// The status code of the response is 404:
	require.True(t, response.IsStatusCode(404))

	// The body contains the error response of the webserver:
	body, err := response.GetBodyAsString()
	require.NoError(t, err)
	require.EqualValues(t, "404 page not found\n", body)
}
