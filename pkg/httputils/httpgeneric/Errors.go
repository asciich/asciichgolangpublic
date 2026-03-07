package httpgeneric

import "errors"

var ErrWebServerAlreadyRunning = errors.New("web server already running")
var ErrUnexpectedStatusCode = errors.New("unexpected status code")
var ErrChecksumMismatch = errors.New("checksum mismatch")

func IsErrorWebServeralreadyRunning(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrUnexpectedStatusCode)
}

func IsErrorUnexpectedStatusCode(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrUnexpectedStatusCode)
}

func IsErrorChecksumMismatch(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrChecksumMismatch)
}
