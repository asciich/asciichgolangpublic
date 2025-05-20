package httputils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/checksums"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getClientByImplementationName(implementationName string) (client Client) {
	if implementationName == "nativeClient" {
		return NewNativeClient()
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
				const verbose bool = true
				const port int = 9123

				testServer, err := GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(verbose)

				err = testServer.StartInBackground(verbose)
				require.NoError(t, err)

				var client Client = getClientByImplementationName(tt.implementationName)
				var response Response
				response, err = client.SendRequest(
					&RequestOptions{
						Url:     "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())),
						Verbose: verbose,
						Method:  tt.method,
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
				const verbose bool = true
				const port int = 9123

				testServer, err := GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(verbose)

				err = testServer.StartInBackground(verbose)
				require.NoError(t, err)

				var client Client = getClientByImplementationName(tt.implementationName)
				responseBody, err := client.SendRequestAndGetBodyAsString(
					&RequestOptions{
						Url:     "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())),
						Verbose: verbose,
						Method:  tt.method,
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
				const verbose bool = true
				const port int = 9123

				testServer, err := GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(verbose)

				err = testServer.StartInBackground(verbose)
				require.NoError(t, err)

				tempFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				defer tempFile.MustDelete(verbose)

				const expectedOutput = "hello world\n"

				var client Client = getClientByImplementationName(tt.implementationName)
				_, err = client.DownloadAsFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:     "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Verbose: verbose,
							Method:  tt.method,
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

				testServer, err := GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(verbose)

				err = testServer.StartInBackground(verbose)
				require.NoError(t, err)

				tempFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				defer tempFile.MustDelete(verbose)

				const expectedOutput = "hello world\n"

				var client Client = getClientByImplementationName(tt.implementationName)
				downloadedFile, err := client.DownloadAsFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Method: tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  checksums.GetSha256SumFromString(expectedOutput),
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)
				defer downloadedFile.MustDelete(verbose)

				require.EqualValues(t, expectedOutput, downloadedFile.MustReadAsString())

				downloadedFile, err = client.DownloadAsFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Method: tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  checksums.GetSha256SumFromString(expectedOutput),
						Verbose:    verbose,
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
				const verbose bool = true
				const port int = 9123

				testServer, err := GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(verbose)

				err = testServer.StartInBackground(verbose)
				require.NoError(t, err)

				var client Client = getClientByImplementationName(tt.implementationName)
				downloadedFile, err := client.DownloadAsTemporaryFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:     "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
							Verbose: verbose,
							Method:  tt.method,
						},
					},
				)
				require.NoError(t, err)
				defer downloadedFile.MustDelete(verbose)

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
				const verbose bool = true
				const port int = 9123

				testServer, err := GetTestWebServer(port)
				require.NoError(t, err)
				defer testServer.Stop(verbose)

				err = testServer.StartInBackground(verbose)
				require.NoError(t, err)

				var client Client = getClientByImplementationName(tt.implementationName)
				output, err := client.SendRequestAndRunYqQueryAgainstBody(
					&RequestOptions{
						Url:     "http://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/example1.yaml",
						Verbose: verbose,
						Method:  tt.method,
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

				const verbose bool = true
				const port int = 9123

				testServer := mustutils.Must(GetTlsTestWebServer(ctx, port))
				defer testServer.Stop(verbose)

				err := testServer.StartInBackground(verbose)
				require.NoError(t, err)

				var client Client = getClientByImplementationName(tt.implementationName)
				output, err := client.SendRequestAndGetBodyAsString(
					&RequestOptions{
						Url:               "https://localhost:" + strconv.Itoa(mustutils.Must(testServer.GetPort())) + "/hello_world.txt",
						Verbose:           verbose,
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
