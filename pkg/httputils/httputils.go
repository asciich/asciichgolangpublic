package httputils

import (
	"context"
	"net/http"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpnativeclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SendRequest(ctx context.Context, requestOptions *httpoptions.RequestOptions) (httputilsinterfaces.Response, error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	return httpnativeclientoo.NewNativeClient().SendRequest(ctx, requestOptions)
}

func SendRequestAndGetBodyAsBytes(ctx context.Context, requestOptions *httpoptions.RequestOptions) (response []byte, err error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	return httpnativeclientoo.NewNativeClient().SendRequestAndGetBodyAsBytes(ctx, requestOptions)
}

func SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *httpoptions.RequestOptions) (response string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	return httpnativeclientoo.NewNativeClient().SendRequestAndGetBodyAsString(ctx, requestOptions)
}

func DownloadAsFile(ctx context.Context, options *httpoptions.DownloadAsFileOptions) (downloadedFile filesinterfaces.File, err error) {
	return httpnativeclientoo.NewNativeClient().DownloadAsFile(ctx, options)
}

func DownloadAsTemporaryFile(ctx context.Context, options *httpoptions.DownloadAsTemporaryFileOptions) (downloadedFile filesinterfaces.File, err error) {
	return httpnativeclientoo.NewNativeClient().DownloadAsTemporaryFile(ctx, options)
}

func WaitUntilStatusCodeOK(ctx context.Context, url string, timeout time.Duration) error {
	if url == "" {
		return tracederrors.TracedErrorEmptyString("url")
	}

	logging.LogInfoByCtxf(ctx, "Wait until %s returns HTTP status code OK started.", url)

	startTime := time.Now()
	deadline := startTime.Add(timeout)

	for {
		if time.Now().After(deadline) {
			return tracederrors.TracedErrorf("Timeout after '%v' waiting for HTTP status code OK from '%s'.", timeout, url)
		}

		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				logging.LogInfoByCtxf(ctx, "%s returned 200 OK.", url)
				break
			}
		}

		elapsed := time.Since(startTime).Round(time.Second)
		remaining := time.Until(deadline).Round(time.Second)
		logging.LogInfoByCtxf(ctx, "Still waiting for HTTP status code OK from '%s'. Elapsed: %v, remaining: %v.", url, elapsed, remaining)

		select {
		case <-ctx.Done():
			return tracederrors.TracedErrorf("Context cancelled while waiting for HTTP status code OK from '%s': %w", url, ctx.Err())
		case <-time.After(1 * time.Second):
		}
	}

	logging.LogInfoByCtxf(ctx, "Wait until %s returns HTTP status code OK finished.", url)

	return nil
}
