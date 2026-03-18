package ioutils

import "io"

// A simple implementation of the io.ReadCloser able to be initialized with lambda functions:
type ReadCloser struct {
	CloseFunc func() error
	ReadFunc  func(p []byte) (n int, err error)
}

func (r *ReadCloser) Close() error {
	if r.CloseFunc == nil {
		return nil
	}

	return r.CloseFunc()
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	if r.ReadFunc == nil {
		return 0, io.EOF
	}

	return r.ReadFunc(p)
}
