package httputilsinterfaces

type Response interface {
	CheckStatusCode(expectedStatusCodes []int) error 
	GetBodyAsBytes() (body []byte, err error)
	GetBodyAsString() (body string, err error)
	IsStatusCode(expectedStatusCode int) bool
	IsStatusCode200Ok() bool
	SetBody(body []byte) (err error)
	SetStatusCode(statusCode int) (err error)
	RunJqQueryAgainstBody(query string) (result string, err error)
	RunYqQueryAgainstBody(query string) (result string, err error)
}
