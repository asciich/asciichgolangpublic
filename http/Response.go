package http

type Response interface {
	GetBodyAsString() (body string, err error)
	IsStatusCodeOk() (isStatusCodeOK bool, err error)
	MustGetBodyAsString() (body string)
	MustIsStatusCodeOk() (isStatusCodeOK bool)
	SetBody(body []byte) (err error)
	SetStatusCode(statusCode int) (err error)
}
