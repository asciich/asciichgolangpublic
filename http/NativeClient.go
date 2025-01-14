package http

import (
	"io"
	"net/http"

	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

// HTTP client written using native go http implementation.
type NativeClient struct {
}

// Get the HTTP client written using native go http implementation.
//
// This is the default client to use when sending request from your running machine.
func GetNativeClient() (client Client) {
	return NewNativeClient()
}

func NewNativeClient() (n *NativeClient) {
	return new(NativeClient)
}

func (c *NativeClient) SendRequest(requestOptions *RequestOptions) (response Response, err error) {
	if requestOptions == nil {
		return nil, errors.TracedErrorNil("requestOptions")
	}

	url, err := requestOptions.GetUrl()
	if err != nil {
		return nil, err
	}

	method, err := requestOptions.GetMethod()
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	nativeResponse, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer nativeResponse.Body.Close()

	response = NewGenericResponse()
	body, err := io.ReadAll(nativeResponse.Body)
	if err != nil {
		return nil, errors.TracedErrorf("Unable to read body as bytes: %w", err)
	}

	err = response.SetBody(body)
	if err != nil {
		return nil, err
	}

	err = response.SetStatusCode(nativeResponse.StatusCode)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (n *NativeClient) MustSendRequest(requestOptions *RequestOptions) (response Response) {
	response, err := n.SendRequest(requestOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return response
}
