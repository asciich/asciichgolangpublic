package httputils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SendRequest(ctx context.Context, requestOptions *httpoptions.RequestOptions) (httputilsinterfaces.Response, error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	return GetNativeClient().SendRequest(ctx, requestOptions)
}

func SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *httpoptions.RequestOptions) (response string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	return GetNativeClient().SendRequestAndGetBodyAsString(ctx, requestOptions)
}

func DownloadAsFile(ctx context.Context, options *httpoptions.DownloadAsFileOptions) (downloadedFile filesinterfaces.File, err error) {
	return GetNativeClient().DownloadAsFile(ctx, options)
}

type progressEveryNBytes struct{}

func WithDownloadProgressEveryNMBytes(ctx context.Context, nMBytes int) context.Context {
	return WithDownloadProgressEveryNkBytes(ctx, 1024*nMBytes)
}

func WithDownloadProgressEveryNkBytes(ctx context.Context, nkBytes int) context.Context {
	return WithDownloadProgressEveryNBytes(ctx, 1024*nkBytes)
}

func WithDownloadProgressEveryNBytes(ctx context.Context, nBytes int) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if nBytes < 0 {
		nBytes = 0
	}

	return context.WithValue(ctx, progressEveryNBytes{}, nBytes)
}

func GetProgressEveryNBytes(ctx context.Context) int {
	if ctx == nil {
		return 0
	}

	val := ctx.Value(progressEveryNBytes{})
	if val == nil {
		return 0
	}

	ret, ok := val.(int)
	if !ok {
		return 0
	}

	return ret
}
