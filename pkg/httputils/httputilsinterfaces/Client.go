package httputilsinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
)

type Client interface {
	DownloadAsFile(ctx context.Context, downloadOptions *httpoptions.DownloadAsFileOptions) (downloadedFile filesinterfaces.File, err error)
	DownloadAsTemporaryFile(ctx context.Context, downloadOptions *httpoptions.DownloadAsTemporaryFileOptions) (downloadedFile filesinterfaces.File, err error)
	SendRequest(ctx context.Context, requestOptions *httpoptions.RequestOptions) (response Response, err error)
	SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *httpoptions.RequestOptions) (responseBody string, err error)
	SendRequestAndRunYqQueryAgainstBody(ctx context.Context, requestOptions *httpoptions.RequestOptions, query string) (result string, err error)

	// Set the base URL to use for all requests.
	// Setting the base URL on the client simplyfies the request as it has to be set only once.
	SetBaseUrl(baseUrl string) error

	// Set the port of the webserver to use for all requests.
	// Setting the port on the client simplyfies the requests as it has to be set only once.
	SetPort(port int) error

	// Set the basic auth authentication to use for all requests.
	SetBasicAuth(*httpoptions.BasicAuth) error
}
