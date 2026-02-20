package httputils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpnativeclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SendRequest(ctx context.Context, requestOptions *httpoptions.RequestOptions) (httputilsinterfaces.Response, error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	return httpnativeclientoo.NewNativeClient().SendRequest(ctx, requestOptions)
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
