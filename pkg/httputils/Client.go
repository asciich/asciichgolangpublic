package httputils

import (
	"context"
	"errors"

	"github.com/asciich/asciichgolangpublic/files"
)

var ErrUnexpectedStatusCode = errors.New("unexpected status code")

type Client interface {
	DownloadAsFile(ctx context.Context, downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	DownloadAsTemporaryFile(ctx context.Context, downloadOptions *DownloadAsFileOptions) (downloadedFile files.File, err error)
	SendRequest(ctx context.Context, requestOptions *RequestOptions) (response Response, err error)
	SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *RequestOptions) (responseBody string, err error)
	SendRequestAndRunYqQueryAgainstBody(ctx context.Context, requestOptions *RequestOptions, query string) (result string, err error)
}
