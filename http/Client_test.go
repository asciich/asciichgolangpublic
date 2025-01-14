package http

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/logging"
)

func getClientByImplementationName(implementationName string) (client Client) {
	if implementationName == "nativeClient" {
		return NewNativeClient()
	}

	logging.LogFatalWithTracef(
		"Unknwon implmentation name '%s'",
		implementationName,
	)

	return nil
}

func Test_Client_GetRequest_RootPage_PortInUrl(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeClient"},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
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
						Url:     "http://localhost:" + strconv.Itoa(port),
						Verbose: verbose,
						Method:  "GET",
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
