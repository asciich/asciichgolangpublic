package httpgeneric

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type HttpHeader struct {
	StatusCode int
}

func (h *HttpHeader) GetStatusCode() (int, error) {
	if h.StatusCode <= 0 {
		return 0, tracederrors.TracedErrorf("Status code not set")
	}

	return h.StatusCode, nil
}
