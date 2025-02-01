package httputils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/logging"
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
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

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

				assert.True(response.MustIsStatusCodeOk())
				assert.Contains(
					response.MustGetBodyAsString(),
					"TestWebServer",
				)
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
				assert := assert.New(t)

				const verbose bool = true
				const port int = 9123

				testServer := MustGetTestWebServer(port)
				defer testServer.Stop(verbose)

				testServer.MustStartInBackground(verbose)

				tempFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				defer tempFile.MustDelete(verbose)

				var client Client = getClientByImplementationName(tt.implementationName)
				downloadedFile := client.MustDownloadAsFile(
					&DownloadAsFileOptions{
						RequestOptions: &RequestOptions{
							Url:     "http://localhost:" + strconv.Itoa(testServer.MustGetPort()) + "/hello_world.txt",
							Verbose: verbose,
							Method:  tt.method,
						},
						OutputPath: tempFile.MustGetPath(),
					},
				)
				defer downloadedFile.MustDelete(verbose)

				assert.Contains(
					"hello world\n",
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
				assert := assert.New(t)

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

				assert.Contains(
					"hello world\n",
					downloadedFile.MustReadAsString(),
				)
			},
		)
	}
}
