package httputils

import (
	"github.com/asciich/asciichgolangpublic/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// This is the generic response type.
// It can also be seen as the default response to use.
type GenericResponse struct {
	body       []byte
	statusCode int
}

func NewGenericResponse() (g *GenericResponse) {
	return new(GenericResponse)
}

func (g *GenericResponse) RunYqQueryAgainstBody(query string) (result string, err error) {
	if query == "" {
		return "", tracederrors.TracedErrorEmptyString("query")
	}

	body, err := g.GetBodyAsString()
	if err != nil {
		return "", err
	}

	return yamlutils.RunYqQueryAginstYamlStringAsString(body, query)
}

func (g *GenericResponse) RunJqQueryAgainstBody(query string) (result string, err error) {
	if query == "" {
		return "", tracederrors.TracedErrorEmptyString("query")
	}

	body, err := g.GetBodyAsString()
	if err != nil {
		return "", err
	}

	return jsonutils.RunJqAgainstJsonStringAsString(body, query)
}

func (g *GenericResponse) GetBody() (body []byte, err error) {
	if g.body == nil {
		return nil, tracederrors.TracedErrorf("body not set")
	}

	return g.body, nil
}

func (g *GenericResponse) GetBodyAsString() (body string, err error) {
	bodyBytes, err := g.GetBody()
	if err != nil {
		return "", err
	}

	return string(bodyBytes), err
}

func (g *GenericResponse) GetStatusCode() (statusCode int, err error) {
	if g.statusCode <= 0 {
		return -1, tracederrors.TracedError("statusCode not set")
	}

	return g.statusCode, nil
}

func (g *GenericResponse) CheckStatusCode(expectedStatusCode int) error {
	if !g.IsStatusCode(expectedStatusCode) {
		return tracederrors.TracedErrorf("%w: %d does not match expected status code %d", ErrUnexpectedStatusCode, g.statusCode, expectedStatusCode)
	}

	return nil
}

func (g *GenericResponse) IsStatusCode(expectedStatusCode int) bool {
	statusCode, err := g.GetStatusCode()
	if err != nil {
		return false
	}

	return statusCode == expectedStatusCode
}

func (g *GenericResponse) IsStatusCode200Ok() bool {
	return g.IsStatusCode(STATUS_CODE_OK)
}

func (g *GenericResponse) SetBody(body []byte) (err error) {
	if body == nil {
		return tracederrors.TracedErrorf("body is nil")
	}

	if len(body) <= 0 {
		return tracederrors.TracedErrorf("body has no elements")
	}

	g.body = body

	return nil
}

func (g *GenericResponse) SetStatusCode(statusCode int) (err error) {
	if statusCode <= 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for statusCode", statusCode)
	}

	g.statusCode = statusCode

	return nil
}
