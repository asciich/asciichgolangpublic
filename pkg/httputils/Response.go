package httputils

type Response interface {
	GetBodyAsString() (body string, err error)
	IsStatusCodeOk() (isStatusCodeOK bool, err error)
	SetBody(body []byte) (err error)
	SetStatusCode(statusCode int) (err error)
	RunJqQueryAgainstBody(query string) (result string, err error)
	RunYqQueryAgainstBody(query string) (result string, err error)
}
