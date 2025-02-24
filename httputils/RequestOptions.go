package httputils

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type RequestOptions struct {
	// Url to request
	Url string

	// Port to use.
	// Overrides the port defined in URL if specified.
	Port int

	// Request method like GET, POST...
	Method string

	// Skip TLS validation
	SkipTLSvalidation bool

	// Enable verbose output
	Verbose bool
}

func NewRequestOptions() (r *RequestOptions) {
	return new(RequestOptions)
}

func (r *RequestOptions) GetDeepCopy() (copy *RequestOptions) {
	copy = new(RequestOptions)

	*copy = *r

	return copy
}

func (r *RequestOptions) GetMethod() (method string, err error) {
	if r.Method == "" {
		return "", tracederrors.TracedErrorf("Method not set")
	}

	return strings.ToUpper(r.Method), nil
}

func (r *RequestOptions) GetMethodOrDefault() (method string, err error) {
	if r.IsMethodSet() {
		return r.GetMethod()
	}

	return "GET", err
}

func (r *RequestOptions) GetPort() (port int, err error) {
	if r.Port <= 0 {
		return -1, tracederrors.TracedError("Port not set")
	}

	return r.Port, nil
}

func (r *RequestOptions) GetSkipTLSvalidation() (skipTLSvalidation bool) {

	return r.SkipTLSvalidation
}

func (r *RequestOptions) GetUrl() (url string, err error) {
	if r.Url == "" {
		return "", tracederrors.TracedErrorf("Url not set")
	}

	return r.Url, nil
}

func (r *RequestOptions) GetVerbose() (verbose bool) {

	return r.Verbose
}

func (r *RequestOptions) IsMethodSet() (isSet bool) {
	return r.Method != ""
}

func (r *RequestOptions) MustGetMethod() (method string) {
	method, err := r.GetMethod()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return method
}

func (r *RequestOptions) MustGetMethodOrDefault() (method string) {
	method, err := r.GetMethodOrDefault()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return method
}

func (r *RequestOptions) MustGetPort() (port int) {
	port, err := r.GetPort()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return port
}

func (r *RequestOptions) MustGetUrl() (url string) {
	url, err := r.GetUrl()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return url
}

func (r *RequestOptions) MustSetMethod(method string) {
	err := r.SetMethod(method)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RequestOptions) MustSetPort(port int) {
	err := r.SetPort(port)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RequestOptions) MustSetUrl(url string) {
	err := r.SetUrl(url)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RequestOptions) SetMethod(method string) (err error) {
	if method == "" {
		return tracederrors.TracedErrorf("method is empty string")
	}

	r.Method = method

	return nil
}

func (r *RequestOptions) SetPort(port int) (err error) {
	if port <= 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for port", port)
	}

	r.Port = port

	return nil
}

func (r *RequestOptions) SetSkipTLSvalidation(skipTLSvalidation bool) {
	r.SkipTLSvalidation = skipTLSvalidation
}

func (r *RequestOptions) SetUrl(url string) (err error) {
	if url == "" {
		return tracederrors.TracedErrorf("url is empty string")
	}

	r.Url = url

	return nil
}

func (r *RequestOptions) SetVerbose(verbose bool) {
	r.Verbose = verbose
}
