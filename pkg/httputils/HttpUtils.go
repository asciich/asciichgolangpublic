package httputils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func SendRequest(ctx context.Context, requestOptions *RequestOptions) (Response, error) {
	if requestOptions == nil {
		return nil, tracederrors.TracedErrorNil("requestOptions")
	}

	return GetNativeClient().SendRequest(ctx, requestOptions)
}

func SendRequestAndGetBodyAsString(ctx context.Context, requestOptions *RequestOptions) (response string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	return GetNativeClient().SendRequestAndGetBodyAsString(ctx, requestOptions)
}
