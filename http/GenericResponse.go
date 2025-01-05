package http

import "github.com/asciich/asciichgolangpublic"

// This is the generic response type.
// It can also be seen as the default response to use.
type GenericResponse struct {
	body       []byte
	statusCode int
}

func NewGenericResponse() (g *GenericResponse) {
	return new(GenericResponse)
}

func (g *GenericResponse) GetBody() (body []byte, err error) {
	if g.body == nil {
		return nil, asciichgolangpublic.TracedErrorf("body not set")
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
		return -1, asciichgolangpublic.TracedError("statusCode not set")
	}

	return g.statusCode, nil
}

func (g *GenericResponse) IsStatusCodeOk() (isStatusCodeOk bool, err error) {
	statusCode, err := g.GetStatusCode()
	if err != nil {
		return false, err
	}

	return statusCode == STATUS_CODE_OK, nil
}

func (g *GenericResponse) MustGetBody() (body []byte) {
	body, err := g.GetBody()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return body
}

func (g *GenericResponse) MustGetBodyAsString() (body string) {
	body, err := g.GetBodyAsString()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return body
}

func (g *GenericResponse) MustGetStatusCode() (statusCode int) {
	statusCode, err := g.GetStatusCode()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return statusCode
}

func (g *GenericResponse) MustIsStatusCodeOk() (isStatusCodeOk bool) {
	isStatusCodeOk, err := g.IsStatusCodeOk()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isStatusCodeOk
}

func (g *GenericResponse) MustSetBody(body []byte) {
	err := g.SetBody(body)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (g *GenericResponse) MustSetStatusCode(statusCode int) {
	err := g.SetStatusCode(statusCode)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (g *GenericResponse) SetBody(body []byte) (err error) {
	if body == nil {
		return asciichgolangpublic.TracedErrorf("body is nil")
	}

	if len(body) <= 0 {
		return asciichgolangpublic.TracedErrorf("body has no elements")
	}

	g.body = body

	return nil
}

func (g *GenericResponse) SetStatusCode(statusCode int) (err error) {
	if statusCode <= 0 {
		return asciichgolangpublic.TracedErrorf("Invalid value '%d' for statusCode", statusCode)
	}

	g.statusCode = statusCode

	return nil
}