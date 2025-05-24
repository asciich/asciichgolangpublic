package httputils_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/checksums"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getClientByImplementationName(implementationName string) (client httputils.Client) {
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

				var client httputils.Client = getClientByImplementationName(tt.implementationName)
				var response httputils.Response
				response, err = client.SendRequest(
					ctx,
					&httputils.RequestOptions{
						Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())),
						Method: tt.method,
					},
				)
				require.NoError(t, err)

				require.True(t, mustutils.Must(response.IsStatusCodeOk()))
				require.Contains(t, mustutils.Must(response.GetBodyAsString()), "TestWebServer")
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

				var client httputils.Client = getClientByImplementationName(tt.implementationName)
				responseBody, err := client.SendRequestAndGetBodyAsString(
					ctx,
					&httputils.RequestOptions{
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

				var client httputils.Client = getClientByImplementationName(tt.implementationName)
				_, err = client.DownloadAsFile(
					ctx,
					&httputils.DownloadAsFileOptions{
						RequestOptions: &httputils.RequestOptions{
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

				var client httputils.Client = getClientByImplementationName(tt.implementationName)
				downloadedFile, err := client.DownloadAsFile(
					ctx,
					&httputils.DownloadAsFileOptions{
						RequestOptions: &httputils.RequestOptions{
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
					&httputils.DownloadAsFileOptions{
						RequestOptions: &httputils.RequestOptions{
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

				var client httputils.Client = getClientByImplementationName(tt.implementationName)
				downloadedFile, err := client.DownloadAsTemporaryFile(
					ctx,
					&httputils.DownloadAsFileOptions{
						RequestOptions: &httputils.RequestOptions{
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

				var client httputils.Client = getClientByImplementationName(tt.implementationName)
				output, err := client.SendRequestAndRunYqQueryAgainstBody(
					ctx,
					&httputils.RequestOptions{
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

				var client httputils.Client = getClientByImplementationName(tt.implementationName)
				output, err := client.SendRequestAndGetBodyAsString(
					ctx,
					&httputils.RequestOptions{
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
