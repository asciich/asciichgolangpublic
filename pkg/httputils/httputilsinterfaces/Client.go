package httputilsinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsparameteroptions"
)

type Client interface {
	DownloadAsFile(ctx context.Context, downloadOptions *httputilsparameteroptions.DownloadAsFileOptions) (downloadedFile files.File, err error)
	DownloadAsTemporaryFile(ctx context.Context, downloadOptions *httputilsparameteroptions.DownloadAsFileOptions) (downloadedFile files.File, err error)
	SendRequest(ctx context.Context, requestOptions *httputilsparameteroptions.RequestOptions) (response Response, err error)
	SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *httputilsparameteroptions.RequestOptions) (responseBody string, err error)
	SendRequestAndRunYqQueryAgainstBody(ctx context.Context, requestOptions *httputilsparameteroptions.RequestOptions, query string) (result string, err error)
}
