package httpgeneric

import "errors"

var ErrWebServerAlreadyRunning = errors.New("web server already running")
var ErrUnexpectedStatusCode = errors.New("unexpected status code")
