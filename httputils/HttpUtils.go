package httputils

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func MustSendRequestAndGetBodyAsString(requestOptions *RequestOptions) (response string) {
	response, err := SendRequestAndGetBodyAsString(requestOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return response
}

func SendRequestAndGetBodyAsString(requestOptions *RequestOptions) (response string, err error) {
	if requestOptions == nil {
		return "", tracederrors.TracedErrorNil("requestOptions")
	}

	return GetNativeClient().SendRequestAndGetBodyAsString(requestOptions)
}
