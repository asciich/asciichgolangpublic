package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsparameteroptions"
)

// Example how to perform a get request and work with the received json data.
func Test_Example_GetJsonDataAndRunJq(t *testing.T) {
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

	// This is the example document we will download from the webserver:
	expectedJson := `{"hello": "world"}`

	// Use a get request to receive the YAML data:
	response, err := httputils.SendRequest(
		ctx,
		&httputilsparameteroptions.RequestOptions{
			// Add the URL to request here:
			Url: "http://localhost:9123/example1.json",

			// There is no need to specify the default method "GET":
			// Method: "GET",
		},
	)
	require.NoError(t, err)

	// Validate the response of the webserver matches the expectedYaml.
	// HINT: This is only done to make the example more understandable.
	bodyString, err := response.GetBodyAsString()
	require.NoError(t, err)
	require.EqualValues(t, expectedJson, bodyString)

	// Extract the value behind 'hello' out ouf the yaml response:
	value, err := response.RunJqQueryAgainstBody(".hello")
	require.NoError(t, err)
	require.EqualValues(t, "world", value)
}
