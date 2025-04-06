package httputils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/checksums"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/mustutils"
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
				require := require.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := MustGetTestWebServer(port)
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				var client Client = getClientByImplementationName(tt.implementationName)
				var response Response = client.MustSendRequest(
					&RequestOptions{
						Url:     "http://localhost:" + strconv.Itoa(testServer.MustGetPort()),
						Verbose: verbose,
						Method:  tt.method,
					},
				)

				require.True(response.MustIsStatusCodeOk())
				require.Contains(
					response.MustGetBodyAsString(),
					"TestWebServer",
				)
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
				require := require.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := MustGetTestWebServer(port)
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				var client Client = getClientByImplementationName(tt.implementationName)
				responseBody := client.MustSendRequestAndGetBodyAsString(
					&RequestOptions{
						Url:     "http://localhost:" + strconv.Itoa(testServer.MustGetPort()),
						Verbose: verbose,
						Method:  tt.method,
					},
				)

				require.Contains(
					responseBody,
					"TestWebServer",
				)
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
				require := require.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := MustGetTestWebServer(port)
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				tempFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				defer tempFile.MustDelete(verbose)

				const expectedOutput = "hello world\n"

				var client Client = getClientByImplementationName(tt.implementationName)
				_, err := client.DownloadAsFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:     "http://localhost:" + strconv.Itoa(testServer.MustGetPort()) + "/hello_world.txt",
							Verbose: verbose,
							Method:  tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  "a" + checksums.GetSha256SumFromString(expectedOutput),
					},
				)
				require.Error(err)
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
				require := require.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := MustGetTestWebServer(port)
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				tempFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				defer tempFile.MustDelete(verbose)

				const expectedOutput = "hello world\n"

				var client Client = getClientByImplementationName(tt.implementationName)
				downloadedFile := client.MustDownloadAsFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(testServer.MustGetPort()) + "/hello_world.txt",
							Method: tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  checksums.GetSha256SumFromString(expectedOutput),
						Verbose:    verbose,
					},
				)
				defer downloadedFile.MustDelete(verbose)

				require.EqualValues(
					expectedOutput,
					downloadedFile.MustReadAsString(),
				)

				downloadedFile = client.MustDownloadAsFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:    "http://localhost:" + strconv.Itoa(testServer.MustGetPort()) + "/hello_world.txt",
							Method: tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
						Sha256Sum:  checksums.GetSha256SumFromString(expectedOutput),
						Verbose:    verbose,
					},
				)
				require.EqualValues(
					expectedOutput,
					downloadedFile.MustReadAsString(),
				)
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
				require := require.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := MustGetTestWebServer(port)
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				var client Client = getClientByImplementationName(tt.implementationName)
				downloadedFile := client.MustDownloadAsTemporaryFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:     "http://localhost:" + strconv.Itoa(testServer.MustGetPort()) + "/hello_world.txt",
							Verbose: verbose,
							Method:  tt.method,
						},
					},
				)
				defer downloadedFile.MustDelete(verbose)

				require.Contains(
					"hello world\n",
					downloadedFile.MustReadAsString(),
				)
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
				require := require.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := MustGetTestWebServer(port)
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				var client Client = getClientByImplementationName(tt.implementationName)
				output := client.MustSendRequestAndRunYqQueryAgainstBody(
					&RequestOptions{
						Url:     "http://localhost:" + strconv.Itoa(testServer.MustGetPort()) + "/example1.yaml",
						Verbose: verbose,
						Method:  tt.method,
					},
					".hello",
				)

				require.EqualValues(
					"world",
					output,
				)
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
				require := require.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := mustutils.Must(GetTlsTestWebServer(ctx, port))
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				var client Client = getClientByImplementationName(tt.implementationName)
				output := client.MustSendRequestAndGetBodyAsString(
					&RequestOptions{
						Url:               "https://localhost:" + strconv.Itoa(testServer.MustGetPort()) + "/hello_world.txt",
						Verbose:           verbose,
						Method:            tt.method,
						SkipTLSvalidation: true,
					},
				)

				require.Contains(output, "hello world")
			},
		)
	}
}
