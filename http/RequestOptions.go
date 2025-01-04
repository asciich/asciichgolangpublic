package http

import (
	"strings"

	"github.com/asciich/asciichgolangpublic"
)

type RequestOptions struct {
	// Url to request
	Url string

	// Port to use.
	// Overrides the port defined in URL if specified.
	Port int

	// Request method like GET, POST...
	Method string

	// Enable verbose output
	Verbose bool
}

func NewRequestOptions() (r *RequestOptions) {
	return new(RequestOptions)
}

func (r *RequestOptions) GetMethod() (method string, err error) {
	if r.Method == "" {
		return "", asciichgolangpublic.TracedErrorf("Method not set")
	}

	return strings.ToUpper(r.Method), nil
}

func (r *RequestOptions) GetPort() (port int, err error) {
	if r.Port <= 0 {
		return -1, asciichgolangpublic.TracedError("Port not set")
	}

	return r.Port, nil
}

func (r *RequestOptions) GetUrl() (url string, err error) {
	if r.Url == "" {
		return "", asciichgolangpublic.TracedErrorf("Url not set")
	}

	return r.Url, nil
}

func (r *RequestOptions) GetVerbose() (verbose bool) {

	return r.Verbose
}

func (r *RequestOptions) MustGetMethod() (method string) {
	method, err := r.GetMethod()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return method
}

func (r *RequestOptions) MustGetPort() (port int) {
	port, err := r.GetPort()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return port
}

func (r *RequestOptions) MustGetUrl() (url string) {
	url, err := r.GetUrl()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return url
}

func (r *RequestOptions) MustSetMethod(method string) {
	err := r.SetMethod(method)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (r *RequestOptions) MustSetPort(port int) {
	err := r.SetPort(port)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (r *RequestOptions) MustSetUrl(url string) {
	err := r.SetUrl(url)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (r *RequestOptions) SetMethod(method string) (err error) {
	if method == "" {
		return asciichgolangpublic.TracedErrorf("method is empty string")
	}

	r.Method = method

	return nil
}

func (r *RequestOptions) SetPort(port int) (err error) {
	if port <= 0 {
		return asciichgolangpublic.TracedErrorf("Invalid value '%d' for port", port)
	}

	r.Port = port

	return nil
}

func (r *RequestOptions) SetUrl(url string) (err error) {
	if url == "" {
		return asciichgolangpublic.TracedErrorf("url is empty string")
	}

	r.Url = url

	return nil
}

func (r *RequestOptions) SetVerbose(verbose bool) {
	r.Verbose = verbose
}
