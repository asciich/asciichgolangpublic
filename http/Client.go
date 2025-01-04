package http

type Client interface {
	SendRequest(requestOptions *RequestOptions) (response Response, err error)
	MustSendRequest(requestOptions *RequestOptions) (response Response)
}
