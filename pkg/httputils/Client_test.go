package httputils_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/checksums"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getClientByImplementationName(implementationName string) (client httputilsinterfaces.Client) {
	if implementationName == "nativeClient" {
		return httputils.NewNativeClient()
	}

	logging.LogFatalWithTracef(
		"Unknown implmentation name '%s'",
		implementationName,
	)

	return nil
}

func TestClient_GetRequest_RootPage_PortInUrl(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
		{"nativeClient", "Get"},
		{"nativeClient", "GET"},
		{"nativeClient", "GeT"},
		{"nativeClient", ""},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const port int = 9123
				ctx := getCtx()

				testServer, err := httputils.GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(ctx)

				err = testServer.StartInBackground(ctx)
				require.NoError(t, err)

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				var response httputilsinterfaces.Response
				response, err = client.SendRequest(
					ctx,
					&httputilsparameteroptions.RequestOptions{
						Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())),
						Method: tt.method,
					},
				)
				require.NoError(t, err)

				require.True(t, response.IsStatusCode200Ok())
			},
		)
	}
}

func TestClient_GetRequestBodyAsString_RootPage_PortInUrl(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
		{"nativeClient", "Get"},
		{"nativeClient", "GET"},
		{"nativeClient", "GeT"},
		{"nativeClient", ""},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const port int = 9123

				testServer, err := httputils.GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(ctx)

				err = testServer.StartInBackground(ctx)
				require.NoError(t, err)

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				responseBody, err := client.SendRequestAndGetBodyAsString(
					ctx,
					&httputilsparameteroptions.RequestOptions{
						Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())),
						Method: tt.method,
					},
				)
				require.NoError(t, err)
				require.Contains(t, responseBody, "TestWebServer")
			},
		)
	}
}

func TestClient_GetRequest_404_PortInUrl(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
		{"nativeClient", "Get"},
		{"nativeClient", "GET"},
		{"nativeClient", "GeT"},
		{"nativeClient", ""},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const port int = 9123

				testServer, err := httputils.GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(ctx)

				err = testServer.StartInBackground(ctx)
				require.NoError(t, err)

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				response, err := client.SendRequest(
					ctx,
					&httputilsparameteroptions.RequestOptions{
						Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/this-page-does-not-exist",
						Method: tt.method,
					},
				)
				require.Error(t, err)
				require.ErrorIs(t, err, httputilsimplementationindependend.ErrUnexpectedStatusCode)

				require.NotNil(t, response)
				require.False(t, response.IsStatusCode200Ok())
				require.True(t, response.IsStatusCode(404))
			},
		)
	}
}

func TestClient_DownloadAsFile_ChecksumMismatch(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const port int = 9123

				testServer, err := httputils.GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(ctx)

				err = testServer.StartInBackground(ctx)
				require.NoError(t, err)

				tempFile := tempfiles.MustCreateEmptyTemporaryFile(contextutils.GetVerboseFromContext(ctx))
				defer tempFile.MustDelete(contextutils.GetVerboseFromContext(ctx))

				const expectedOutput = "hello world\n"

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				_, err = client.DownloadAsFile(
					ctx,
					&httputilsparameteroptions.DownloadAsFileOptions{
						RequestOptions: &httputilsparameteroptions.RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Method: tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  "a" + checksums.GetSha256SumFromString(expectedOutput),
					},
				)
				require.Error(t, err)
			},
		)
	}
}

func TestClient_DownloadAsFile(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				const port int = 9123
				ctx := getCtx()

				testServer, err := httputils.GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(ctx)

				err = testServer.StartInBackground(ctx)
				require.NoError(t, err)

				tempFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				defer tempFile.MustDelete(verbose)

				const expectedOutput = "hello world\n"

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				downloadedFile, err := client.DownloadAsFile(
					ctx,
					&httputilsparameteroptions.DownloadAsFileOptions{
						RequestOptions: &httputilsparameteroptions.RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Method: tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  checksums.GetSha256SumFromString(expectedOutput),
					},
				)
				require.NoError(t, err)
				defer downloadedFile.MustDelete(verbose)

				require.EqualValues(t, expectedOutput, downloadedFile.MustReadAsString())

				downloadedFile, err = client.DownloadAsFile(
					ctx,
					&httputilsparameteroptions.DownloadAsFileOptions{
						RequestOptions: &httputilsparameteroptions.RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Method: tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  checksums.GetSha256SumFromString(expectedOutput),
					},
				)
				require.NoError(t, err)
				require.EqualValues(t, expectedOutput, downloadedFile.MustReadAsString())
			},
		)
	}
}

func TestClient_DownloadAsTempraryFile(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const port int = 9123

				testServer, err := httputils.GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(ctx)

				err = testServer.StartInBackground(ctx)
				require.NoError(t, err)

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				downloadedFile, err := client.DownloadAsTemporaryFile(
					ctx,
					&httputilsparameteroptions.DownloadAsFileOptions{
						RequestOptions: &httputilsparameteroptions.RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Method: tt.method,
						},
					},
				)
				require.NoError(t, err)
				defer downloadedFile.MustDelete(contextutils.GetVerboseFromContext(ctx))

				require.Contains(t, "hello world\n", downloadedFile.MustReadAsString())
			},
		)
	}
}

func TestClient_GetRequestAndRunYqQuery(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const port int = 9123

				testServer, err := httputils.GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(ctx)

				err = testServer.StartInBackground(ctx)
				require.NoError(t, err)

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				output, err := client.SendRequestAndRunYqQueryAgainstBody(
					ctx,
					&httputilsparameteroptions.RequestOptions{
						Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/example1.yaml",
						Method: tt.method,
					},
					".hello",
				)
				require.NoError(t, err)

				require.EqualValues(t, "world", output)
			},
		)
	}
}

func TestClient_GetRequestUsingTls_insecure(t *testing.T) {
	tests := []struct {
		implementationName string
		method             string
	}{
		{"nativeClient", "get"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const port int = 9123

				testServer := mustutils.Must(httputils.GetTlsTestWebServer(ctx, port))
				defer testServer.Stop(ctx)

				err := testServer.StartInBackground(ctx)
				require.NoError(t, err)

				var client httputilsinterfaces.Client = getClientByImplementationName(tt.implementationName)
				output, err := client.SendRequestAndGetBodyAsString(
					ctx,
					&httputilsparameteroptions.RequestOptions{
						Url:               "https://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
						Method:            tt.method,
						SkipTLSvalidation: true,
					},
				)
				require.NoError(t, err)
				require.Contains(t, output, "hello world")
			},
		)
	}
}
