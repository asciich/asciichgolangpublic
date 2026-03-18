package ioutils

import "errors"

var ErrWriteFunctionNotSet = errors.New("write function not set")

// A simple implementation of the io.ReadCloser able to be initialized with lambda functions:
type WriteCloser struct {
	CloseFunc func() error
	WriteFunc func(p []byte) (n int, err error)
}

func (w *WriteCloser) Close() error {
	if w.CloseFunc == nil {
		return nil
	}

	return w.CloseFunc()
}

func (w *WriteCloser) Write(p []byte) (n int, err error) {
	if w.WriteFunc == nil {
		return 0, ErrWriteFunctionNotSet
	}

	return w.WriteFunc(p)
}
